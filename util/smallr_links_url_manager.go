package util

import (
	"fmt"

	"github.com/adwitiyaio/arka/logger"
)

const smallrLinksApiUrlKey = "SMALLR_LINKS_API_URL"
const smallrLinksApiKey = "SMALLR_LINKS_API_KEY"

type smallrLinksProvider struct {
	BaseUrl        string
	ApiKey         string
	CreateLinkPath string
}

type smallrLinksShortenUrlDto struct {
	Target string `json:"target"`
}

type smallLink struct {
	Shortcode  string `json:"shortcode"`
	Target     string `json:"target"`
	ExpiryDate string `json:"expiry_date"`
}

func (mus *multiUrlManager) initializeSmallrLinks() {
	smallrLinksBaseUrl := mus.sm.GetValueForKey(smallrLinksApiUrlKey)
	smallrLinksApiKey := mus.sm.GetValueForKey(smallrLinksApiKey)
	mus.slp = &smallrLinksProvider{
		BaseUrl:        smallrLinksBaseUrl,
		ApiKey:         smallrLinksApiKey,
		CreateLinkPath: "/",
	}
}

func (mus *multiUrlManager) shortenWithSmallrLinks(url string) (string, error) {
	reqBody := smallrLinksShortenUrlDto{
		Target: url,
	}

	path := fmt.Sprintf("%s%s", mus.slp.BaseUrl, mus.slp.CreateLinkPath)
	var response smallLink
	resp, err := mus.client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("x-api-key", mus.slp.ApiKey).
		SetBody(reqBody).
		SetResult(&response).
		Post(path)

	if err != nil {
		logger.Log.Error().Err(err).Msg("failed to shorten url with smallr links")
		return "", err
	}

	logger.Log.Debug().Msgf("smallr links response -> %s", string(resp.Body()))
	return fmt.Sprintf("%s/%s", mus.slp.BaseUrl, response.Shortcode), nil
}
