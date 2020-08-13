package dynamodb

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/nik/platform-image-service/logger"
	"github.com/nik/platform-image-service/pkg/domain/model"
)

const tableName = "platform_image_db.api_metadata"

type DynamoDBAPIMetadataRepo struct {
	dynamodb *dynamodb.DynamoDB
}

func NewDynamoDBAPIMetadataRepo(svc *dynamodb.DynamoDB) *DynamoDBAPIMetadataRepo {
	repo := &DynamoDBAPIMetadataRepo{
		dynamodb: svc,
	}

	return repo
}

//get retrieves api metadata for a provided tenant identifier and apiname
func (repo *DynamoDBAPIMetadataRepo) Get(tenantID string, apiName string) (*model.APIMetadata, error) {
	logger := logger.GetInstance()
	logger.Sugar().Infof("Retrieving data from api_metada for tenant_id - %s and api_name - %s", tenantID, apiName)
	// Query to retrieve metadata from aws_metadata table
	result, err := repo.dynamodb.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"tenant_id": {
				S: aws.String(tenantID),
			},
			"api_name": {
				S: aws.String(apiName),
			},
		},
	})

	if err != nil {
		logger.Sugar().Infof("Retrieving data from api_metada failed with error %s", err.Error())
		return nil, err
	}

	//unmarshal the object into target struct of api_metadata
	item := &model.APIMetadata{}
	err = dynamodbattribute.UnmarshalMap(result.Item, &item)
	if err != nil {
		//attribute mismatch
		logger.Sugar().Infof("Unmarshaling of the result failed with error %s", err.Error())
	}

	return item, nil
}
