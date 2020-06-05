package service

import (
	"encoding/json"
	"fmt"
	"github.com/nik/platform-image-service/pkg/domain/model"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

const apiName = "google_image_search_api"

type ImageService struct {
	apiService *APIService
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

	instance.search(&searchRequest)
	return nil, error
}

func (instace *ImageService) search(request *model.ImageRequest) (response *model.ImageSearchResponse) {
	key, apiUrl, searchEngineId, error := instace.apiService.GetAPIKeyUrlAndSearchEngineID(request.TenantID, apiName)
	if error != nil {
		return nil
	}

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
