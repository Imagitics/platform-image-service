package messaging

type MessagingServiceInterface interface {
	Publish(stream *string, partitionKey string, event string) (*MessageResponse, error)
}
