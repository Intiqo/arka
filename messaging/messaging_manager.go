package messaging

import (
	"github.com/adwitiyaio/arka/dependency"
	"github.com/adwitiyaio/arka/secrets"
)

const DependencyMessagingManager = "messaging_manager"

const ProviderFirebase = "firebase"
const ProviderOneSignal = "onesignal"

type Message struct {
	Title    string
	Body     string
	ImageUrl string
	Data     map[string]interface{}
	Tokens   []string
	Channel  string
}

type Manager interface {
	// SendNotification ... Send a push notification
	//
	// Currently supported providers are:
	//
	// - firebase
	//
	// - onesignal
	//
	// Returns the raw response, a list of failed tokens and any error
	SendNotificationWithProvider(message Message, provider string) (interface{}, []string, error)
}

// Bootstrap ... Bootstrap the messaging manager
func Bootstrap() {
	dm := dependency.GetManager()
	mm := &multiManager{
		sm: dm.Get(secrets.DependencySecretsManager).(secrets.Manager),
	}
	mm.initialize()
	dm.Register(DependencyMessagingManager, mm)
}
