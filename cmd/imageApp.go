package main

import (
	"github.com/nik/platform-image-service/pkg/domain/repository"
	"github.com/nik/platform-image-service/pkg/domain/service"
	"github.com/nik/platform-image-service/pkg/infra"
)

func main() {
	//create cassandra connection instance
	conn := &cassandra.CassandraConn{
		Hosts:       []string{"172.18.0.2"},
		Port:        "9042",
		User:        "cassandra",
		Password:    "cassandra",
		Consistency: "Quorum",
		Keyspace:    "platform_image_db",
	}

	//create repoinstance
	repoInstance := repository.NewCassandraAPIMetadataRepo(conn)
	apiServiceInstance := service.NewAPIService(repoInstance)
	imageSearch := service.NewImageService(apiServiceInstance)
	imageSearch.Search("TEST01", "cars")
}
