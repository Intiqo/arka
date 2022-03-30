package sms

import (
	"encoding/base64"
	"fmt"
	"os"

	"github.com/adwitiyaio/arka/logger"
)

const clickSendUsernameKey = "CLICKSEND_USERNAME"
const clickSendApiKey = "CLICKSEND_API_KEY"

type clickSendProvider struct {
	baseUrl     string
	sendSmsPath string
	authHeader  string
}

type clickSendMsg struct {
	To     string `json:"to"`
	Source string `json:"source"`
	Body   string `json:"body"`
}

type clickSendSmsBody struct {
	Messages []clickSendMsg `json:"messages"`
}

type clickSendResponseBody struct {
	HTTPCode     int    `json:"http_code"`
	ResponseCode string `json:"response_code"`
	ResponseMsg  string `json:"response_msg"`
	Data         struct {
		TotalPrice  float64 `json:"total_price"`
		TotalCount  int     `json:"total_count"`
		QueuedCount int     `json:"queued_count"`
		Messages    []struct {
			Direction    string      `json:"direction"`
			Date         int         `json:"date"`
			To           string      `json:"to"`
			Body         string      `json:"body"`
			From         string      `json:"from"`
			Schedule     interface{} `json:"schedule"`
			MessageID    string      `json:"message_id"`
			MessageParts int         `json:"message_parts"`
			MessagePrice interface{} `json:"message_price"`
			CustomString string      `json:"custom_string"`
			UserID       int         `json:"user_id"`
			SubaccountID int         `json:"subaccount_id"`
			Country      string      `json:"country"`
			Carrier      string      `json:"carrier"`
			Status       string      `json:"status"`
		} `json:"messages"`
		Currency struct {
			CurrencyNameShort string `json:"currency_name_short"`
			CurrencyPrefixD   string `json:"currency_prefix_d"`
			CurrencyPrefixC   string `json:"currency_prefix_c"`
			CurrencyNameLong  string `json:"currency_name_long"`
		} `json:"currency"`
	} `json:"data"`
}

func (msm *multiSmsManager) initializeClickSend() {
	username := msm.cm.GetValueForKey(clickSendUsernameKey)
	apiKey := msm.cm.GetValueForKey(clickSendApiKey)

	auth := fmt.Sprintf("%s:%s", username, apiKey)
	encAuth := base64.StdEncoding.EncodeToString([]byte(auth))
	ah := fmt.Sprintf("Basic %s", encAuth)

	msm.csc = &clickSendProvider{
		baseUrl:     "https://rest.clicksend.com/v3",
		sendSmsPath: "/sms/send",
		authHeader:  ah,
	}
}

func (msm multiSmsManager) sendSmsViaClickSend(options Options) (interface{}, error) {
	// Create a series of messages based on the number of recipients
	msgs := make([]clickSendMsg, 0)
	for _, rec := range options.Recipients {
		msgs = append(msgs, clickSendMsg{
			To:     rec,
			Source: "app",
			Body:   options.Message,
		})
	}

	reqBody := clickSendSmsBody{Messages: msgs}
	url := fmt.Sprintf("%s%s", msm.csc.baseUrl, msm.csc.sendSmsPath)

	if os.Getenv("CI") != "true" {
		return msm.dispatchClickSend(reqBody, url)
	}
	return nil, nil
}

func (msm multiSmsManager) dispatchClickSend(reqBody clickSendSmsBody, url string) (interface{}, error) {
	var response clickSendResponseBody
	resp, err := msm.client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", msm.csc.authHeader).
		SetBody(reqBody).
		SetResult(&response).
		Post(url)

	if err != nil {
		logger.Log.Error().Err(err).Msg("failed to send sms via click send")
		return nil, err
	}
	return resp.Body(), nil
}
