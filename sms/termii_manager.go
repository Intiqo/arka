package sms

import (
	"errors"
	"os"

	termii "github.com/Uchencho/go-termii"

	"github.com/adwitiyaio/arka/logger"
	"github.com/adwitiyaio/arka/secrets"
)

const termiiApiKey = "TERMII_API_KEY"

type termiiManager struct {
	sm secrets.Manager

	apiKey string
	client termii.Client
}

func (tm *termiiManager) initialize() {
	tm.client = termii.NewClient()
	tm.apiKey = tm.sm.GetValueForKey(termiiApiKey)
}

func (tm termiiManager) SendSms(options Options) (interface{}, error) {
	if len(options.Recipients) < 1 {
		return nil, errors.New("no recipients")
	}

	responses := make([]termii.AutoGeneratedMessageResponse, 0)
	for _, recipient := range options.Recipients {
		req := termii.AutoGeneratedMessageRequest{
			To:     recipient,
			Sms:    options.Message,
			APIKey: tm.apiKey,
		}
		if os.Getenv("CI") != "true" {
			resp := tm.dispatchTermii(req, recipient)
			responses = append(responses, resp)
		}
	}

	return responses, nil
}

func (tm termiiManager) dispatchTermii(
	req termii.AutoGeneratedMessageRequest,
	recipient string,
) termii.AutoGeneratedMessageResponse {
	resp, err := tm.client.SendAutoGeneratedMessage(req)
	if err != nil {
		logger.Log.Error().Err(err).Msgf("failed to send SMS to recipient %s", recipient)
	}
	return resp
}
