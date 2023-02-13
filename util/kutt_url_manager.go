package util

import (
	"fmt"

	"github.com/adwitiyaio/arka/logger"
)

const kuttApiUrlKey = "KUTT_API_URL"
const kuttApiKey = "KUTT_API_KEY"

type kuttProvider struct {
	BaseUrl        string
	ApiKey         string
	CreateLinkPath string
}

type kuttParams struct {
	Target      string `json:"target"`
	Description string `json:"description"`
	ExpireIn    string `json:"expire_in"`
	Password    string `json:"password"`
	CustomUrl   string `json:"customurl"`
	Reuse       bool   `json:"reuse"`
	Domain      string `json:"domain"`
}

type kuttUrl struct {
	ID          string      `json:"id"`
	Address     string      `json:"address"`
	Description interface{} `json:"description"`
	Banned      bool        `json:"banned"`
	Password    bool        `json:"password"`
	ExpireIn    interface{} `json:"expire_in"`
	Target      string      `json:"target"`
	VisitCount  int64       `json:"visit_count"`
	CreatedAt   string      `json:"created_at"`
	UpdatedAt   string      `json:"updated_at"`
	Link        string      `json:"link"`
}

func (mus *multiUrlManager) initializeKutt() {
	kuttBaseUrl := mus.sm.GetValueForKey(kuttApiUrlKey)
	kuttApiKey := mus.sm.GetValueForKey(kuttApiKey)
	mus.kp = &kuttProvider{
		BaseUrl:        kuttBaseUrl,
		ApiKey:         kuttApiKey,
		CreateLinkPath: "/links",
	}
}

func (mus *multiUrlManager) shortenWithKutt(url string) (string, error) {
	reqBody := kuttParams{
		Target: url,
	}

	path := fmt.Sprintf("%s/api/v2%s", mus.kp.BaseUrl, mus.kp.CreateLinkPath)
	var response kuttUrl
	resp, err := mus.client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("X-API-Key", mus.kp.ApiKey).
		SetBody(reqBody).
		SetResult(&response).
		Post(path)

	if err != nil {
		logger.Log.Error().Err(err).Msg("failed to shorten url with kutt")
		return "", err
	}

	logger.Log.Debug().Msgf("http response -> %s", string(resp.Body()))
	return response.Link, nil
}
