package service

import "github.com/nik/platform-image-service/pkg/domain/model"

type ImageServiceInterface interface {
	Search(request model.ImageRequest) model.ImageSearchResponse
}

type APIServiceInterface interface {
	GetAPIKeyAndSearchEngineID(tenantID string, apiName string) (string, string, string, error)
}
