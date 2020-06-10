package repository

import (
	"errors"
	"github.com/gocql/gocql"
	"github.com/nik/platform-image-service/pkg/domain/model"
	"github.com/nik/platform-image-service/pkg/infra/cassandra"
)

type CassandraImageStoreMetadataRepo struct {
	session *gocql.Session
}

func NewCassandraImageStoreMetadataRepo(conn *cassandra.CassandraConn) *CassandraImageStoreMetadataRepo {
	session := conn.InitSession()
	repo := &CassandraImageStoreMetadataRepo{
		session: session,
	}

	return repo
}

// GetMetadataByTenantID retrieves aws metadata for a provided tenant identifier.
// This metadata includes secret key, access key and preferred_region.
func (repo *CassandraImageStoreMetadataRepo) Get(tenantID string, searchTerm string, searchAlias string) (*model.ImageStoreData, error) {
	imageStoreData := &model.ImageStoreData{}

	// Query to retrieve metadata from aws_metadata table
	selectQuery := "select image_count, store_type, image_store_by_title from image_metadata where tenant_id = ? and searchTerm = ? and searchAlias = ?"
	iter := repo.session.Query(selectQuery, tenantID, searchTerm, searchAlias).Iter()

	if iter == nil {
		//some problem with cassandra
		return nil, errors.New("Invalid request")
	} else if iter.NumRows() == 0 {
		//no data indicates that t
		return imageStoreData, nil
	} else {
		// Scan and store relevant attributes into struct
		m := map[string]interface{}{}
		iter.MapScan(m)
		imageStoreData := &model.ImageStoreData{
			TenantId:           tenantID,
			ImageCount:         m["image_count"].(int),
			Searchterm:         searchTerm,
			SearchTermAlias:    searchAlias,
			StoreType:          m["store_type"].(string),
			StoreUrlByImageURL: m["image_store_by_title"].(map[string]string),
		}
		return imageStoreData, nil
	}
}
