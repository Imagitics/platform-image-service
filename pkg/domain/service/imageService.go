package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/nik/platform-image-service/pkg/domain/model"
	"github.com/nik/platform-image-service/pkg/domain/repository"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

const apiName = "google_image_search_api"
const maxImages = 1000
const batchSize = 10

type ImageService struct {
	apiService *APIService
	repo       repository.CassandraImageStoreMetadataRepo
}

func NewImageService(apiServiceInstance *APIService) (imageSearchInstance *ImageService) {
	return &ImageService{
		apiService: apiServiceInstance,
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

	fmt.Println(data)

	return &data
}

func (instance *ImageService) SearchAndCollectImages(tenantID string, searchTerm string, searchTermAlias string) error {
	//retrieve api metadata for this tenant
	apiKey, apiUrl, searchEngineId, error := instance.apiService.GetAPIKeyUrlAndSearchEngineID(tenantID, apiName)

	if error != nil {
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
	if imageCount < maxImages {
		return nil
	} else {
		//search images and start storing images
		storeUrlByImageTitle := map[string]string{}
		for counter := imageCount; counter < maxImages; counter = counter + 10 {
			searchRequest := &model.ImageRequest{TenantID: tenantID,
				APIKey:         apiKey,
				SearchEngineID: searchEngineId,
				SearchTerm:     searchTerm,
				APIUrl:         apiUrl,
				Start:          imageCount,
				End:            imageCount + 10,
				IncludeFace:    false,
			}

			//invoke search api with searchRequest
			imageRes := instance.search(searchRequest)
			for counter := 0; counter < len(imageRes.Items); counter++ {
				imageTitle := imageRes.Items[counter].Title
				storeUrlByImageTitle[imageTitle] = imageRes.Items[counter].Link
			}
		}
	}
	return nil
}
