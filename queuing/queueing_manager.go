package queuing

import (
	"errors"

	"github.com/adwitiyaio/arka/cloud"
	"github.com/adwitiyaio/arka/dependency"
	"github.com/adwitiyaio/arka/logger"
	"github.com/adwitiyaio/arka/secrets"
)

const DependencyQueuingManager = "queueing_manager"

const ProviderSQS = "sqs"

type Options struct {
	QueueName string
}

type SendOptions struct {
	Options
	GroupId string
	Data    interface{}
}

type ReceiveOptions struct {
	Options
	DelayTimeout     int64
	WaitTimeSeconds  int64
	NumberOfMessages int64
}

type MessageResponse struct {
	MessageId string
	Data      interface{}
	Receipt   string
}

type ReceiveResponse struct {
	Messages []MessageResponse
}

type DeleteOptions struct {
	Options
	MessageHandle string
}

type Manager interface {
	// SendMessage ... Send message to a queue
	SendMessage(options SendOptions) error
	// ReceiveMessage ... Receive messages from a queue
	ReceiveMessage(options ReceiveOptions) (ReceiveResponse, error)
	// DeleteMessage ... Delete messages from a queue
	DeleteMessage(options DeleteOptions) error
}

// Bootstrap ... Bootstraps the cloud manager
func Bootstrap(provider string) {
	dm := dependency.GetManager()
	var mm interface{}
	switch provider {
	case ProviderSQS:
		mm = &sqsManager{
			sm:  dm.Get(secrets.DependencySecretsManager).(secrets.Manager),
			clm: dm.Get(cloud.DependencyCloudManager).(cloud.Manager),
		}
		mm.(*sqsManager).initialize()
	default:
		err := errors.New("queuing provider not implemented")
		logger.Log.Fatal().Err(err).Str("provider", provider)
	}
	dm.Register(DependencyQueuingManager, mm)
}
