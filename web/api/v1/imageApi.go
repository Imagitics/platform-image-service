package v1

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/nik/platform-image-service/logger"
	"github.com/nik/platform-image-service/pkg/domain/service"
	"github.com/nik/platform-image-service/web/rest/model"
	"go.uber.org/zap"
	"net/http"
)

const FileSizeLimitError = "multipart: NextPart: http: request body too large"
const NoSuchFileError = "http: no such file"

type Api struct {
	router             *mux.Router
	imageSearchService *service.ImageService
}

func NewApi(router *mux.Router, imageService *service.ImageService) *Api {
	s3Handler := &Api{
		router:             router,
		imageSearchService: imageService,
	}

	return s3Handler
}

// It is search handler for searching the images based on search term
func (api *Api) search(w http.ResponseWriter, r *http.Request) {
	imageSearchRequest, errorCode, err := validateAndRetriveSearchRequest(r)
	if err != nil {
		respondWithJSON(w, errorCode, fmt.Sprintf(err.Error()))
	} else {
		searchResults, err := api.imageSearchService.Search(imageSearchRequest.TenantID, imageSearchRequest.SearchTerm)
		if err != nil {
			respondWithJSON(w, http.StatusInternalServerError, fmt.Sprintf("No results found for search term %s", imageSearchRequest.SearchTerm))
		} else if searchResults == nil || len(searchResults.Items) == 0 {
			respondWithJSON(w, http.StatusOK, fmt.Sprintf("No results found for search term %s", imageSearchRequest.SearchTerm))
		} else {
			urls := []model.ImageURL{}
			for counter := 0; counter < len(searchResults.Items); counter++ {
				url := model.ImageURL{
					Name: searchResults.Items[counter].Title,
					Url:  searchResults.Items[counter].Link,
				}

				urls = append(urls, url)
			}

			imageSearchResponse := model.SearchAPIImageResponse{TotalResults: len(searchResults.Items),
				SearchUrls: urls,
			}
			respondWithJSON(w, http.StatusOK, imageSearchResponse)
		}
	}
}

//collectImages collects images for the search term and upload to store type
func (api *Api) collectImages(w http.ResponseWriter, r *http.Request) {
	logger := logger.GetInstance()
	imageSearchRequest, errorCode, err := validateAndRetriveSearchRequest(r)
	if err != nil {
		logger.Info("Responding with ", zap.String("error code", err.Error()))
		respondWithJSON(w, errorCode, fmt.Sprintf(err.Error()))
	} else {
		//err := api.imageSearchService.PublishCollectImageEvent(imageSearchRequest)
		err := api.imageSearchService.SearchAndCollectImages(imageSearchRequest.TenantID, imageSearchRequest.SearchTerm, imageSearchRequest.SearchAlias)
		if err != nil {
			logger.Info("Responding with ", zap.String("error code", err.Error()))
			respondWithJSON(w, http.StatusInternalServerError, fmt.Sprintf("No results found for search term %s", imageSearchRequest.SearchTerm))
		} else {
			logger.Info("Responding with success")
			respondWithJSON(w, http.StatusOK, fmt.Sprintf("Request successfully accepted"))
		}
	}
}

//validateSearchTerm validates the search term as per conformed rules
func validateSearchTerm(searchTerm string) error {
	if searchTerm == "" {
		return errors.New("Search term can not be empty")
	}

	return nil
}

//validateSearchTerm validates the search term as per conformed rules
func validateTenant(tenantID string) error {
	if tenantID == "" {
		return errors.New("Invalid tenant")
	}

	return nil
}

//validateSearchTerm validates the search term as per conformed rules
func retrieveIncludeFace(includeFace string) (bool, error) {
	if includeFace == "" {
		return false, nil
	} else if includeFace == "true" {
		return true, nil
	} else if includeFace == "false" {
		return false, nil
	}

	return false, errors.New("Invalid include_face argument")
}

//validateAndRetriveSearchRequest validates the search request
//It validates and retrieves tenant identifier and search term
func validateAndRetriveSearchRequest(r *http.Request) (*model.SearchAPIImageRequest, int, error) {
	// validate tenant
	vars := mux.Vars(r)
	tenantId := vars["tenant_id"]
	if err := validateTenant(tenantId); err != nil {
		return nil, 401, err
	}

	searchTerm := r.URL.Query().Get("search_term")
	if err := validateSearchTerm(searchTerm); err != nil {
		return nil, 400, err
	}

	searchAlias := r.URL.Query().Get("search_alias")
	appName := r.URL.Query().Get("app_name")
	searchAlias = retrieveSearchTermAlias(searchAlias, searchTerm, appName)

	includeFace := false
	if value, err := retrieveIncludeFace(r.URL.Query().Get("include_face")); err != nil {
		return nil, 400, err
	} else {
		includeFace = value
	}

	imageRequest := &model.SearchAPIImageRequest{
		TenantID:    tenantId,
		SearchTerm:  searchTerm,
		IncludeFace: includeFace,
		SearchAlias: searchAlias,
	}

	return imageRequest, 0, nil
}

//decorate error response
func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

//decorate success response
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func (a *Api) InitializeRoutes() {
	//a.router.HandleFunc("/s3/images/{id:[0-9]+}", a.getProducts).Methods("GET")
	a.router.HandleFunc("/{tenant_id}/images/search", a.search).Methods("GET")
	a.router.HandleFunc("/{tenant_id}/images/collect", a.collectImages).Methods("GET")
	//a.Router.HandleFunc("/product/{id:[0-9]+}", a.getProduct).Methods("GET")
	//a.Router.HandleFunc("/product/{id:[0-9]+}", a.updateProduct).Methods("PUT")
	//a.Router.HandleFunc("/product/{id:[0-9]+}", a.deleteProduct).Methods("DELETE")
}

//retrieveSearchTermAlias validates the search term alias
func retrieveSearchTermAlias(searchTermAlias string, searchTerm string, appName string) string {
	if searchTermAlias == "" {
		//check whether searchTermAlias is empty
		searchTermAlias = searchTerm
	}
	if appName != "" {
		searchTermAlias = appName + "/" + searchTermAlias
	}

	return searchTermAlias
}
