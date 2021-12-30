package config

import (
	"errors"
	"github.com/adwitiyaio/arka/dependency"
	"github.com/adwitiyaio/arka/logger"
)

const DependencyConfigManager = "config_manager"

const ProviderEnvironment = "ENVIRONMENT"

type Manager interface {
	// GetValueForKey ... Gets the value for a configuration key
	GetValueForKey(key string) string
}

// Bootstrap ... Bootstraps the config manager
func Bootstrap(provider string, configPath string) {
	c := dependency.GetManager()
	var cm interface{}
	switch provider {
	case ProviderEnvironment:
		cm = &environmentConfigManager{}
		cm.(*environmentConfigManager).initialize(configPath)
	default:
		err := errors.New("config provider not implemented")
		logger.Log.Fatal().Err(err).Str("provider", provider)
	}
	c.Register(DependencyConfigManager, cm)
}
