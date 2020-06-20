package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/nik/platform-image-service/config"
	"github.com/nik/platform-image-service/logger"
	"github.com/nik/platform-image-service/pkg/domain/model"
	"github.com/nik/platform-image-service/pkg/domain/repository"
	"github.com/nik/platform-image-service/utility"
	"go.uber.org/zap"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

const apiName = "google_image_search_api"
const maxImages = 100
const batchSize = 10
const rootBucketName = "imagiticstest01"

type ImageService struct {
	apiService *APIService
	repo       repository.ImageStore
	config     *config.ConfigModel
}

func NewImageService(apiServiceInstance *APIService, imageStoreMetadataRepo repository.ImageStore, config *config.ConfigModel) (imageSearchInstance *ImageService) {
	return &ImageService{
		apiService: apiServiceInstance,
		repo:       imageStoreMetadataRepo,
		config:     config,
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

func (instance *ImageService) SearchAndCollectImages(tenantID string, searchTerm string, searchTermAlias string) error {
	//retrieve api metadata for this tenant
	apiKey, apiUrl, searchEngineId, error := instance.apiService.GetAPIKeyUrlAndSearchEngineID(tenantID, apiName)
	logger := logger.GetInstance()

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

		for counter := 1; counter < maxImages; counter = counter + batchSize {
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
				instance.uploadImageToStore(tenantID, searchTermAlias, imageRes.Items[imageCounter].Link, rootFilePath+strconv.Itoa(imageCount)+fileFormat)
				storeUrlByImageTitle[imageTitle] = imageRes.Items[imageCounter].Link
			}
		}
		os.Remove(rootFilePath)
	}
	return nil
}

func (instance *ImageService) uploadImageToStore(tenantID string, searchTermAlias string, linkUrl string, filePath string) (int, error) {
	logger := logger.GetInstance()

	//create s3 upload request
	s3FileUploadReq := &model.S3UploadRequest{
		Bucket:    rootBucketName + tenantID,
		TenantId:  tenantID,
		Directory: searchTermAlias,
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
		fmt.Println("Error in downloading the file")
	}
	request, err := utility.NewfileUploadRequest(instance.config.Platform_S3_URL, extraParams, "entity", filePath)
	if err != nil {
		log.Fatal(err)
	}
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		logger.Info("File upload operation to s3 failed with", zap.String("error", err.Error()))
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

	return http.StatusCreated, nil
}
