package cache

import (
	"errors"

	"github.com/adwitiyaio/arka/dependency"
	"github.com/adwitiyaio/arka/logger"
	"github.com/adwitiyaio/arka/secrets"
)

const DependencyCacheManager = "cache_manager"

const ProviderRedis = "REDIS"
const ProviderLocal = "LOCAL"

type Manager interface {
	// GetStatus ... Returns the current status of the cache system connection
	GetStatus() string

	// Set ... Sets a new value in the cache
	Set(key string, val interface{}) error

	// Get ... Gets a value from cache
	Get(key string) (string, error)
}

// Bootstrap ... Bootstraps the cache manager
func Bootstrap(provider string) {
	var r interface{}
	c := dependency.GetManager()
	switch provider {
	case ProviderRedis:
		r = &redisCacheManager{
			sm: c.Get(secrets.DependencySecretsManager).(secrets.Manager),
		}
		r.(*redisCacheManager).initialize()
	case ProviderLocal:
		r = &localCacheManager{}
		r.(*localCacheManager).initialize()
	default:
		err := errors.New("cache provider not implemented")
		logger.Log.Fatal().Err(err).Str("provider", provider)
	}
	c.Register(DependencyCacheManager, r)
}
