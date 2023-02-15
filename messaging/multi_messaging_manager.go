package messaging

import (
	"errors"

	"github.com/adwitiyaio/arka/secrets"
)

type multiManager struct {
	sm secrets.Manager

	fm  firebaseManager
	osm oneSignalManager
}

func (m *multiManager) initialize() {
	m.fm = firebaseManager{
		sm: m.sm,
	}
	m.fm.initialize()

	m.osm = oneSignalManager{
		sm: m.sm,
	}
	m.osm.initialize()
}

func (m *multiManager) SendNotificationWithProvider(message Message, provider string) (interface{}, []string, error) {
	switch provider {
	case ProviderFirebase:
		return m.fm.SendNotification(message)
	case ProviderOneSignal:
		return m.osm.SendNotification(message)
	default:
		return nil, nil, errors.New("Invalid provider")
	}
}
