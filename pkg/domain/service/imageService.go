package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/nik/platform-image-service/config"
	"github.com/nik/platform-image-service/logger"
	"github.com/nik/platform-image-service/pkg/domain/model"
	"github.com/nik/platform-image-service/pkg/domain/repository"
	"github.com/nik/platform-image-service/pkg/infra/messaging"
	"github.com/nik/platform-image-service/utility"
	model2 "github.com/nik/platform-image-service/web/rest/model"
	"go.uber.org/zap"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

const apiName = "google_image_search_api"
const maxImages = 20
const batchSize = 10
const rootBucketName = "imagitics"

var stream *string

type ImageService struct {
	apiService *APIService
	repo       repository.ImageStore
	config     *config.ConfigModel
	messaging  messaging.MessagingServiceInterface
}

func NewImageService(apiServiceInstance *APIService, imageStoreMetadataRepo repository.ImageStore, messagingInstance messaging.MessagingServiceInterface, config *config.ConfigModel) (imageSearchInstance *ImageService) {
	stream = flag.String("stream", config.Messaging.CollectImageRequestedEventStream, config.Messaging.CollectImageRequestedEventStream)
	return &ImageService{
		apiService: apiServiceInstance,
		repo:       imageStoreMetadataRepo,
		config:     config,
		messaging:  messagingInstance,
	}
}

func (instance *ImageService) Search(tenantID string, searchTerm string) (*model.ImageSearchResponse, error) {
	apiKey, apiUrl, searchEngineId, error := instance.apiService.GetAPIKeyUrlAndSearchEngineID(tenantID, apiName)
	if error != nil {
		return nil, error
	}
	searchRequest := model.ImageRequest{TenantID: tenantID,
		APIKey:         apiKey,
		SearchEngineID: searchEngineId,
		SearchTerm:     searchTerm,
		APIUrl:         apiUrl,
		IncludeFace:    false,
	}

	return instance.search(&searchRequest), nil
}

func (instace *ImageService) search(request *model.ImageRequest) (response *model.ImageSearchResponse) {
	key, apiUrl, searchEngineId := request.APIKey, request.APIUrl, request.SearchEngineID

	//create instance of http request with apiurl
	req, err := http.NewRequest("GET", apiUrl, nil)
	if err != nil {
		os.Exit(1)
	}

	searchEngineId, _ = url.QueryUnescape(searchEngineId)
	key, _ = url.QueryUnescape(key)

	client := &http.Client{}
	q := req.URL.Query()
	q.Add("q", request.SearchTerm)
	q.Add("cx", searchEngineId)
	q.Add("lr", "lang_en")
	q.Add("searchType", "image")
	if request.Start != 0 {
		q.Add("start", strconv.Itoa(request.Start))
	}
	if request.ImagesToSearch != 0 {
		q.Add("num", strconv.Itoa(request.ImagesToSearch))
	}
	//q.Add("imgType", "face")
	q.Add("key", key)
	req.URL.RawQuery = q.Encode()

	fmt.Println(req.URL.String())

	req.Header.Add("Accept", "application/json")

	if err != nil {
		fmt.Println("Errored when sending request to the server")
		return nil
	}
	resp, err := client.Do(req)
	data := model.ImageSearchResponse{}

	if err == nil {
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		_ = json.Unmarshal(body, &data)
	}

	return &data
}

func (instance *ImageService) PublishCollectImageEvent(searchRequst *model2.SearchAPIImageRequest) error {
	logger := logger.GetInstance()
	//get the json message
	searchReqJson, err := json.Marshal(searchRequst)
	if err != nil {
		logger.Sugar().Errorf("Can not  event %s to stream %s failed with error - ", err.Error())
		return err
	}
	//publish the search request
	pubResponse, err := instance.messaging.Publish(stream, searchRequst.TenantID, string(searchReqJson))
	if err != nil {
		logger.Error("Error in producing event to stream", zap.String("error", err.Error()))
		return err
	}

	logger.Sugar().Infof("Event is published to stream %s in partition %s at offset %s", stream, pubResponse.SeqNum, pubResponse.ShardId)
	return nil
}

func (instance *ImageService) SearchAndCollectImages(tenantID string, searchTerm string, searchTermAlias string) error {
	logger := logger.GetInstance()

	//retrieve api metadata for this tenant
	apiKey, apiUrl, searchEngineId, error := instance.apiService.GetAPIKeyUrlAndSearchEngineID(tenantID, apiName)

	if error != nil {
		logger.Info("Responding with ", zap.String("error code", error.Error()))
		return errors.New("Unauthorized request")
	}

	//check whether images for this requests already exist
	//validate whether any delta is remaining
	//by default support only maxImages to store
	prevSearchResults, err := instance.repo.Get(tenantID, searchTerm, searchTermAlias)
	if err != nil {
		return err
	}

	imageCount := prevSearchResults.ImageCount
	if imageCount >= maxImages {
		return nil
	} else {
		//search images and start storing images
		storeUrlByImageTitle := map[string]string{}
		rootFilePath := "/tmp/s3upload/"
		fileFormat := ".jpg"
		os.Mkdir(rootFilePath, 0777)
		totalNumImages := 0

		for counter := 1; totalNumImages < maxImages; counter = counter + batchSize {
			searchRequest := &model.ImageRequest{TenantID: tenantID,
				APIKey:         apiKey,
				SearchEngineID: searchEngineId,
				SearchTerm:     searchTerm,
				APIUrl:         apiUrl,
				Start:          counter,
				IncludeFace:    false,
			}

			//invoke search api with searchRequest
			imageRes := instance.search(searchRequest)

			for imageCounter := 0; imageCounter < len(imageRes.Items); imageCounter++ {
				imageCount = imageCount + 1
				imageTitle := imageRes.Items[imageCounter].Title
				_, err := instance.uploadImageToStore(tenantID, searchTerm, searchTermAlias, imageRes.Items[imageCounter].Link, rootFilePath+strconv.Itoa(imageCount)+fileFormat)
				if err == nil {
					//error implies some problem in uploading the image to s3
					//simply skip this image and move on
					storeUrlByImageTitle[imageTitle] = imageRes.Items[imageCounter].Link
					totalNumImages = totalNumImages + 1
				}
			}
		}

		//perform the checkpointing by updating storage record
		if storeUrlByImageTitle != nil && len(storeUrlByImageTitle) > 0 {
			imageStoreData := &model.ImageStoreData{
				TenantId:           tenantID,
				ImageCount:         totalNumImages,
				Searchterm:         searchTerm,
				SearchTermAlias:    searchTermAlias,
				StoreType:          "aws_s3_storage",
				StoreUrlByImageURL: storeUrlByImageTitle,
			}
			instance.repo.Insert(imageStoreData)
		}

		os.Remove(rootFilePath)
	}
	return nil
}

func (instance *ImageService) uploadImageToStore(tenantID string, searchTerm string, searchTermAlias string, linkUrl string, filePath string) (int, error) {
	logger := logger.GetInstance()

	//create s3 upload request
	s3FileUploadReq := &model.S3UploadRequest{
		Bucket:    rootBucketName + strings.ToLower(tenantID),
		TenantId:  tenantID,
		Directory: searchTermAlias,
		FilePath:  filePath,
	}

	e, err := json.Marshal(s3FileUploadReq)
	if err != nil {
		logger.Info("Responding with ", zap.String("error", err.Error()))
		return http.StatusBadRequest, errors.New("Invalid request")
	}
	//continue to s3 upload operation
	extraParams := map[string]string{
		"request": string(e),
	}

	//download file and upload to s3
	err = utility.DownloadFile(filePath, linkUrl)
	if err != nil {
		return http.StatusInternalServerError, err
	} else {
		request, err := utility.NewfileUploadRequest(instance.config.Platform_S3_URL, extraParams, "entity", filePath)
		if err != nil {
			log.Fatal(err)
		}
		client := &http.Client{}
		resp, err := client.Do(request)
		if err != nil {
			logger.Sugar().Infof("File upload operation to s3 failed with error %s", err.Error())
		} else {
			body := &bytes.Buffer{}
			_, err := body.ReadFrom(resp.Body)
			if err != nil {
				logger.Info("File upload operation to s3 failed with", zap.String("error", err.Error()))
			} else {
				fmt.Println(filePath + "-" + body.String())
			}
			resp.Body.Close()
		}
	}

	return http.StatusCreated, nil
}
