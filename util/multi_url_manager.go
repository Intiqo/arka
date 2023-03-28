package util

import (
	"github.com/go-resty/resty/v2"

	"github.com/adwitiyaio/arka/secrets"
)

type multiUrlManager struct {
	sm       secrets.Manager
	provider string

	client *resty.Client
	fbp    *firebaseDeepLinkProvider
	kp     *kuttProvider
	slp    *smallrLinksProvider
}

func (mus *multiUrlManager) initialize() {
	mus.client = resty.New()
	mus.initializeFirebase()
	switch mus.provider {
	case UrlProviderKutt:
		mus.initializeKutt()
	case UrlProviderSmallrLinks:
		mus.initializeSmallrLinks()
	}
}

func (mus *multiUrlManager) CreateDeepLink(url string) (string, error) {
	return mus.createDeepLinkWithFirebase(url)
}

func (mus *multiUrlManager) Shorten(url string) (string, error) {
	switch mus.provider {
	case UrlProviderKutt:
		return mus.shortenWithKutt(url)
	case UrlProviderSmallrLinks:
		return mus.shortenWithSmallrLinks(url)
	}
	return mus.shortenWithKutt(url)
}
