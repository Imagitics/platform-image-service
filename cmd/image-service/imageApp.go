package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/nik/platform-image-service/logger"
	"github.com/nik/platform-image-service/pkg/domain/repository/dynamodb"
	"github.com/nik/platform-image-service/pkg/domain/service"
	"github.com/nik/platform-image-service/pkg/infra/dynamodb"
	"github.com/nik/platform-image-service/pkg/infra/messaging"
	"github.com/nik/platform-image-service/utility"
	"github.com/nik/platform-image-service/web/api/v1"
	"log"
	"net/http"
	"time"
)

func main() {
	//load configuration
	config, err := utility.LoadConfiguration("/etc/config/config.json")
	if err != nil {
		//halt bootstrapping
		fmt.Println("Error in loading configuration - ", err)
	}

	logger := logger.InitInstance(config.Logger)
	logger.Info("Logger is successfully initialized")
	//instantiate cassandra connection instance
	/*conn := &cassandra.CassandraConn{
		Hosts:       []string{config.Cassandra.Host},
		Port:        config.Cassandra.Port,
		User:        config.Cassandra.User,
		Password:    config.Cassandra.Password,
		Consistency: config.Cassandra.Consistency,
		Keyspace:    config.Cassandra.Keyspace,
		CaPath:      config.Cassandra.SSLCertPath,
	}*/

	//get dynamodb instance
	dynamodbConn := dynamo_db.DynamoDBConn{
		config.Dynamodb.Endpoint,
	}
	dynamodbInstance := dynamodbConn.InitSession(config.AWS)
	logger.Sugar().Infof("dynamodbconnection is initialized")
	//create repoinstance
	repoAPIMetadataInstance := dynamodb.NewDynamoDBAPIMetadataRepo(dynamodbInstance)
	apiServiceInstance := service.NewAPIService(repoAPIMetadataInstance)
	repoImageStoreInstance := dynamodb.NewDynamoDBImageStoreMetadataRepo(dynamodbInstance)
	messaging, err := messaging.NewMessaging(config.Messaging)
	if err == nil {
		imageSearch := service.NewImageService(apiServiceInstance, repoImageStoreInstance, messaging, config)
		//instantiate api object and routes
		router := mux.NewRouter()
		apiInstnace := v1.NewApi(router, imageSearch)
		apiInstnace.InitializeRoutes()

		//create a http server
		srv := &http.Server{
			Addr: ":8080",
			// Good practice: enforce timeouts for servers you create!
			WriteTimeout: 15 * time.Second,
			ReadTimeout:  15 * time.Second,
			Handler:      router,
		}
		fmt.Println("Initializing http server")
		log.Fatal(srv.ListenAndServe())
	}
}
