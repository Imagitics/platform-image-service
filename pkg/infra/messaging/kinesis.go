package messaging

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kinesis"
	"github.com/nik/platform-image-service/config"
	"github.com/nik/platform-image-service/logger"
)

type Messaging struct  {
	kinesis *kinesis.Kinesis
}

type MessageResponse struct {
	SeqNum *string
	ShardId *string
}

//Instantiates new instance of messaging
func NewMessaging(config *config.Messaging) (*Messaging, error) {
	logger := logger.GetInstance()
	s,err := session.NewSession(&aws.Config{Region: aws.String(config.Region)})
	if(err!=nil) {
		logger.Error("Error in creating kinesis session")
		return nil, err
	}
	//create a new instance of kinesis
	kc := kinesis.New(s)
	messaging := &Messaging{
		kinesis: kc,
	}

	return messaging, nil
	}

//Message publisher publishes event to input stream based shared by partition key
func (messaging *Messaging)  Publish(stream *string, partitionKey string, event string) (*MessageResponse, error) {
	//publish the message to stream
	putOutput, err := messaging.kinesis.PutRecord(&kinesis.PutRecordInput{
		Data:         []byte(event),
		StreamName:   stream,
		PartitionKey: aws.String(partitionKey),
	})

	if(err!=nil) {
		return  nil, err
	} else {
		res := &MessageResponse{
			SeqNum:  putOutput.SequenceNumber,
			ShardId: putOutput.ShardId,
		}
		return res, nil
	}
}