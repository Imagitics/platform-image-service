package dynamodb

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/nik/platform-image-service/logger"
	"github.com/nik/platform-image-service/pkg/domain/model"
)

type DynamoDBImageStoreMetadataRepo struct {
	dynamodb *dynamodb.DynamoDB
}

const tableNameImageStore = "platform_image_db.image_store"

func NewDynamoDBImageStoreMetadataRepo(svc *dynamodb.DynamoDB) *DynamoDBImageStoreMetadataRepo {
	repo := &DynamoDBImageStoreMetadataRepo{
		dynamodb: svc,
	}

	return repo
}

// Get retrieves image metadata for a given tenant_id, search_term and search_alias
// This metadata includes secret key, access key and preferred_region.
func (repo *DynamoDBImageStoreMetadataRepo) Get(tenantID string, searchTerm string, searchAlias string) (*model.ImageStoreData, error) {
	logger := logger.GetInstance()
	logger.Sugar().Infof("Retrieving data from platform_image_db.image_store for tenant_id - %s, search_term - %s and search_alias %s", tenantID, searchTerm, searchAlias)
	// Query to retrieve metadata from aws_metadata table
	result, err := repo.dynamodb.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(tableNameImageStore),
		Key: map[string]*dynamodb.AttributeValue{
			"tenant_id": {
				S: aws.String(tenantID),
			},
			"search_term": {
				S: aws.String(searchTerm),
			},
		},
	})

	if err != nil {
		logger.Sugar().Infof("Retrieving data from image_store failed with error %s", err.Error())
		return nil, err
	}

	//unmarshal the object into target struct of image_store
	item := &model.ImageStoreData{}
	err = dynamodbattribute.UnmarshalMap(result.Item, &item)
	if err != nil {
		//attribute mismatch
		logger.Sugar().Infof("Unmarshaling of the result failed with error %s", err.Error())
	}

	return item, nil

}

//Upsert adds metadata into platform_image_db.image_store
func (repo *DynamoDBImageStoreMetadataRepo) Insert(data *model.ImageStoreData) (bool, error) {
	logger := logger.GetInstance()
	item, err := dynamodbattribute.MarshalMap(data)
	if err != nil {
		logger.Sugar().Infof("Unmarshaling of the result failed with error %v", err)
	}

	//map the struct to dynamodb format
insertItem:
	input := &dynamodb.PutItemInput{
		Item:      item,
		TableName: aws.String(tableNameImageStore),
	}
	//put the record into dynamodb
	//check for error. In case of error retry for n times and then log the error
	errorCount := 0
	_, err = repo.dynamodb.PutItem(input)
	if err != nil {
		if errorCount != 3 {
			errorCount++
			goto insertItem
		}
		logger.Sugar().Infof("Unmarshaling of the result failed with error %v", err)
		return false, err
	} else {
		return true, nil
	}
}
