package event

import (
	"errors"

	"github.com/adwitiyaio/arka/cloud"
	"github.com/adwitiyaio/arka/dependency"
	"github.com/adwitiyaio/arka/logger"
)

const DependencyEventManager = "event_manager"

const ProviderLocal = "LOCAL"
const ProviderSns = "SNS"

type Handler func(event string, data interface{})

type Manager interface {
	// Register ... Register an event handler
	Register(event string, handler Handler) error

	// Publish ... Publish an event
	Publish(event string, data interface{}) error
}

func Bootstrap(provider string) {
	dm := dependency.GetManager()
	var m interface{}
	switch provider {
	case ProviderLocal:
		m = &localEventManager{}
		m.(*localEventManager).initialize()
	case ProviderSns:
		m = &snsEventManager{
			clm: dm.Get(cloud.DependencyCloudManager).(cloud.Manager),
		}
		m.(*snsEventManager).initialize()
	default:
		err := errors.New("event manager provider unknown")
		logger.Log.Fatal().Err(err).Str("provider", provider)
	}
	dm.Register(DependencyEventManager, m)
}
