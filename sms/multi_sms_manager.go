package sms

import (
	"fmt"
	"strconv"

	"github.com/go-resty/resty/v2"

	"github.com/adwitiyaio/arka/logger"
	"github.com/adwitiyaio/arka/secrets"
)

type multiSmsManager struct {
	sm     secrets.Manager
	client *resty.Client

	csc *clickSendProvider
	sbc *smsBroadcastProvider
}

func (tm *multiSmsManager) initialize() {
	tm.client = resty.New()
	tm.initializeClickSend()
	tm.initializeSmsBroadcast()
}

func (tm *multiSmsManager) SendSms(options Options) (interface{}, error) {
	const AUCode = "61"
	// We segregate the AU recipients to send SMS via SMS Broadcast
	auRecipients := segregateRecipients(options.Recipients, AUCode)
	if len(auRecipients) > 0 {
		options.Recipients = auRecipients
		return tm.sendSmsViaSmsBroadcast(options)
	}

	// We segregate all other recipients to send SMS via ClickSend
	otherRecipients := segregateRecipients(options.Recipients, "")
	if len(otherRecipients) > 0 {
		options.Recipients = otherRecipients
		return tm.sendSmsViaClickSend(options)
	}

	return nil, nil
}

func segregateRecipients(recipients []string, code string) []string {
	to := make([]string, 0)
	for _, recipient := range recipients {
		num, err := ParsePhoneNumber(recipient)
		if err != nil {
			logger.Log.Debug().Err(err).Msgf("failed to parse phone number, %s", recipient)
			continue
		}
		countryCode := strconv.Itoa(int(*num.CountryCode))
		if countryCode == code || code == "" {
			nn := strconv.Itoa(int(*num.NationalNumber))
			mob := fmt.Sprintf("+%s %s", countryCode, nn)
			to = append(to, mob)
		}
	}
	return to
}
