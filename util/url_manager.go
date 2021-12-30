package util

import (
	"github.com/adwitiyaio/arka/config"
	"github.com/adwitiyaio/arka/dependency"
)

const DependencyUrlManager = "url_manager"

const UrlProviderKutt = "kutt"

type UrlManager interface {
	// Shorten ... Shortens a URL
	Shorten(url string) (string, error)

	// CreateDeepLink ... Creates a deep link for mobile integration
	CreateDeepLink(url string) (string, error)
}

// BootstrapUrlManager ... Bootstraps the URL manager
func BootstrapUrlManager(provider string) {
	dm := dependency.GetManager()
	mus := &multiUrlManager{
		cm:       dm.Get(config.DependencyConfigManager).(config.Manager),
		provider: provider,
	}
	mus.initialize()
	dm.Register(DependencyUrlManager, mus)
}
