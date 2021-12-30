package messaging

import (
	"errors"

	"github.com/adwitiyaio/arka/config"

	"github.com/adwitiyaio/arka/dependency"
	"github.com/adwitiyaio/arka/logger"
)

const DependencyMessagingManager = "messaging_manager"

const ProviderFirebase = "firebase"

type Message struct {
	Title    string
	Body     string
	ImageUrl string
	Data     map[string]string
	Tokens   []string
}

type Manager interface {
	// SendNotification ... Send a push notification
	//
	// @param {Message} message A message object
	// @return {[]string} Returns a list of failed tokens
	SendNotification(message Message) []string
}

// Bootstrap ... Bootstraps the cloud manager
func Bootstrap(provider string) {
	dm := dependency.GetManager()
	var mm interface{}
	switch provider {
	case ProviderFirebase:
		mm = &firebaseManager{
			cm: dm.Get(config.DependencyConfigManager).(config.Manager),
		}
		mm.(*firebaseManager).initialize()
	default:
		err := errors.New("messaging provider not implemented")
		logger.Log.Fatal().Err(err).Str("provider", provider)
	}
	dm.Register(DependencyMessagingManager, mm)
}
