package util

import (
	"github.com/adwitiyaio/arka/config"
	"github.com/go-resty/resty/v2"
)

type multiUrlManager struct {
	cm       config.Manager
	provider string

	client *resty.Client
	fbp    *firebaseDeepLinkProvider
	kp     *kuttProvider
}

func (mus *multiUrlManager) initialize() {
	mus.client = resty.New()
	mus.initializeFirebase()
	mus.initializeKutt()
}

func (mus *multiUrlManager) CreateDeepLink(url string) (string, error) {
	return mus.createDeepLinkWithFirebase(url)
}

func (mus *multiUrlManager) Shorten(url string) (string, error) {
	return mus.shortenWithKutt(url)
}
