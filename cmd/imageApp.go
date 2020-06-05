package main

import (
	"fmt"
	"github.com/nik/platform-image-service/pkg/domain/repository"
	"github.com/nik/platform-image-service/pkg/domain/service"
	"github.com/nik/platform-image-service/pkg/infra/cassandra"
	"github.com/nik/platform-image-service/utility"
)

func main() {

	//load configuration
	config, err := utility.LoadConfiguration("/etc/config/config.json")
	if err != nil {
		//halt bootstrapping
		fmt.Println("Error in loading configuration - ", err)
	}
	//instantiate cassandra connection instance
	conn := &cassandra.CassandraConn{
		Hosts:       []string{config.Cassandra.Host},
		Port:        config.Cassandra.Port,
		User:        config.Cassandra.User,
		Password:    config.Cassandra.Password,
		Consistency: config.Cassandra.Consistency,
		Keyspace:    config.Cassandra.Keyspace,
		CaPath:      config.Cassandra.SSLCertPath,
	}
	//create repoinstance
	repoInstance := repository.NewCassandraAPIMetadataRepo(conn)
	apiServiceInstance := service.NewAPIService(repoInstance)
	//imageMetadataRepo:= repository.NewCassandraAPIMetadataRepo(conn)
	imageSearch := service.NewImageService(apiServiceInstance)
	//imageSearch.Search("TEST01", "cars",0, 10)
	imageSearch.Search("TEST01", "computers")
}
