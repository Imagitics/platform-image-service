package repository

import (
	"errors"
	"github.com/gocql/gocql"
	"github.com/nik/platform-image-service/pkg/domain/model"
	"github.com/nik/platform-image-service/pkg/infra"
)

type CassandraAPIMetadataRepo struct {
	session *gocql.Session
}

func NewCassandraAPIMetadataRepo(conn *cassandra.CassandraConn) *CassandraAPIMetadataRepo {
	conn.Keyspace = "platform_image_db"
	conn.Consistency = "QUORUM"
	session := conn.InitSession()
	repo := &CassandraAPIMetadataRepo{
		session: session,
	}

	return repo
}

// GetMetadataByTenantID retrieves aws metadata for a provided tenant identifier.
// This metadata includes secret key, access key and preferred_region.
func (repo *CassandraAPIMetadataRepo) Get(tenantID string, apiName string) (*model.APIMetadata, error) {
	// Query to retrieve metadata from aws_metadata table
	selectQuery := "select parameters from api_metadata where tenant_id = ? and api_Name = ?"
	iter := repo.session.Query(selectQuery, tenantID, apiName).Iter()
	if iter.NumRows() != 1 {
		//maximum one record is expected as tenant identifier is the unique key
		return nil, errors.New("Bad request")
	}

	// Scan and store relevant attributes into struct
	m := map[string]interface{}{}
	iter.MapScan(m)
	apiMetadataInstance := model.APIMetadata{
		Params: m["parameters"].(map[string]string),
	}

	return &apiMetadataInstance, nil
}
