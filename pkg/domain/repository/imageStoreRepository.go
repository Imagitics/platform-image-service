package repository

import "github.com/nik/platform-image-service/pkg/domain/model"

type ImageStore interface {
	Get(tenantID string, searchTerm string, searchTermAlias string) (*model.ImageStoreData, error)
	Insert(data *model.ImageStoreData) (bool, error)
}
