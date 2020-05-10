package repository

import "github.com/nik/platform-image-service/pkg/domain/model"

type APIRepository interface {
	Get(tenantID string, apiName string) (*model.APIMetadata, error)
}
