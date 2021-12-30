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

func (msm *multiSmsManager) initializeSmsBroadcast() {
	username := msm.cm.GetValueForKey(smsBroadcastUsernameKey)
	password := msm.cm.GetValueForKey(smsBroadcastPasswordKey)

	msm.sbc = &smsBroadcastProvider{
		username:    username,
		password:    password,
		baseUrl:     "https://api.smsbroadcast.com.au",
		sendSmsPath: "/api-adv.php",
	}
}

func (msm multiSmsManager) sendSmsViaSmsBroadcast(options Options) error {
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
	url := fmt.Sprintf("%s%s", msm.sbc.baseUrl, msm.sbc.sendSmsPath)
	if os.Getenv("CI") != "true" {
		return msm.dispatchSmsBroadcast(options, to, url)
	}
	return nil
}

func (msm multiSmsManager) dispatchSmsBroadcast(options Options, to string, url string) error {
	resp, err := msm.client.R().
		SetQueryParams(map[string]string{
			"username": msm.sbc.username,
			"password": msm.sbc.password,
			"to":       to,
			"message":  options.Message,
			"maxsplit": "8",
		}).
		Get(url)

	if err != nil {
		logger.Log.Error().Err(err).Msg("failed to send sms via sms broadcast")
		return err
	}

	logger.Log.Debug().Msgf("http response -> %s", string(resp.Body()))
	return nil
}
