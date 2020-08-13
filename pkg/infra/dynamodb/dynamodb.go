package dynamo_db

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/nik/platform-image-service/config"
	"github.com/nik/platform-image-service/logger"
)

type DynamoDBConn struct {
	Endpoint string
}

func NewDynamoDBConnection(dynamodbConfig config.Dynamodb) *DynamoDBConn {
	dynamoDBConn := &DynamoDBConn{
		Endpoint: dynamodbConfig.Endpoint,
	}

	return dynamoDBConn
}

//InitSession creates a session with dynamodb
func (conn *DynamoDBConn) InitSession(awsConfigModel *config.AWS) *dynamodb.DynamoDB {
	logger := logger.GetInstance()
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(awsConfigModel.Region),
		Credentials: credentials.NewStaticCredentials(awsConfigModel.AccessKey, awsConfigModel.SecretKey, ""),
	})

	if err != nil {
		logger.Sugar().Infof("Instantiation of aws session is failed with %s", err.Error())
	}
	//instantiate the dynamodb service
	svc := dynamodb.New(sess, &aws.Config{Endpoint: aws.String(conn.Endpoint)})
	return svc
}
