package sms

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/adwitiyaio/arka/logger"
)

const smsBroadcastUsernameKey = "SMSBROADCAST_USERNAME"
const smsBroadcastPasswordKey = "SMSBROADCAST_PASSWORD"

type smsBroadcastProvider struct {
	username    string
	password    string
	baseUrl     string
	sendSmsPath string
}

func (tm *multiSmsManager) initializeSmsBroadcast() {
	username := tm.cm.GetValueForKey(smsBroadcastUsernameKey)
	password := tm.cm.GetValueForKey(smsBroadcastPasswordKey)

	tm.sbc = &smsBroadcastProvider{
		username:    username,
		password:    password,
		baseUrl:     "https://api.smsbroadcast.com.au",
		sendSmsPath: "/api-adv.php",
	}
}

func (tm multiSmsManager) sendSmsViaSmsBroadcast(options Options) (interface{}, error) {
	// Create a series of messages based on the number of recipients
	space := regexp.MustCompile(`\s+`)
	to := ""
	for _, rec := range options.Recipients {
		n := strings.ReplaceAll(rec, "+", "")
		vn := space.ReplaceAllString(n, "")
		to = fmt.Sprintf("%s,%s", to, vn)
	}
	to = to[1:]

	// Create a HTTP request and call the API
	url := fmt.Sprintf("%s%s", tm.sbc.baseUrl, tm.sbc.sendSmsPath)
	if os.Getenv("CI") != "true" {
		return tm.dispatchSmsBroadcast(options, to, url)
	}
	return nil, nil
}

func (tm multiSmsManager) dispatchSmsBroadcast(options Options, to string, url string) (interface{}, error) {
	resp, err := tm.client.R().
		SetQueryParams(
			map[string]string{
				"username": tm.sbc.username,
				"password": tm.sbc.password,
				"to":       to,
				"message":  options.Message,
				"maxsplit": "8",
			},
		).
		Get(url)

	if err != nil {
		logger.Log.Error().Err(err).Msg("failed to send sms via sms broadcast")
		return nil, err
	}
	return resp.Body(), nil
}
