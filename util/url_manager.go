package util

import (
	"github.com/adwitiyaio/arka/dependency"
	"github.com/adwitiyaio/arka/secrets"
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
		sm:       dm.Get(secrets.DependencySecretsManager).(secrets.Manager),
		provider: provider,
	}
	mus.initialize()
	dm.Register(DependencyUrlManager, mus)
}
