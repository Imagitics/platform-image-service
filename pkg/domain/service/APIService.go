package service

import "github.com/nik/platform-image-service/pkg/domain/repository"

type APIService struct {
	repo repository.APIRepository
}

func NewAPIService(repoInstance repository.APIRepository) *APIService {
	return &APIService{
		repo: repoInstance,
	}
}

func (apiService *APIService) GetAPIKeyUrlAndSearchEngineID(tenantID string, apiName string) (string, string, string, error) {
	apiMetadata, error := apiService.repo.Get(tenantID, apiName)
	if error != nil {
		return "", "", "", error
	} else {
		return apiMetadata.Params["key"], apiMetadata.Params["url"], apiMetadata.Params["cx"], nil
	}
}
