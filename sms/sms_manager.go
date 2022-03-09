package sms

import (
	"errors"
	"math"
	"unicode/utf8"

	"github.com/nyaruka/phonenumbers"

	"github.com/adwitiyaio/arka/cloud"
	"github.com/adwitiyaio/arka/config"
	"github.com/adwitiyaio/arka/dependency"
	"github.com/adwitiyaio/arka/logger"
)

const DependencySmsManager = "sms_manager"

const singleSmsCharacterCount = 160
const multiSmsCharacterCount = 153
const unicodeSingleSmsCharacterCount = 70
const unicodeMultiSmsCharacterCount = 66

const ProviderMulti = "multi"
const ProviderSns = "sns"

// Options ... Various options to send an SMS.
//
// Recipients is a string array. Recipient should contain the country code as well.
// For example, "+919191092920".
// If any of the mobile number is invalid, it is dropped from the recipient list
//
// A Message can be greater than 160 characters, in which case, the SMS will be split
// into multiple messages
type Options struct {
	Recipients []string
	Message    string
}

// Manager ... SMS Manager that handles sending messages
type Manager interface {
	// SendSms ... Sends an SMS to the recipients.
	//
	// See Options to understand the structure
	SendSms(options Options) error
}

// Bootstrap ... Bootstraps the SMS Manager
func Bootstrap(provider string) {
	dm := dependency.GetManager()
	var smsManager interface{}
	switch provider {
	case ProviderMulti:
		smsManager = &multiSmsManager{
			cm: dm.Get(config.DependencyConfigManager).(config.Manager),
		}
		smsManager.(*multiSmsManager).initialize()
	case ProviderSns:
		smsManager = &snsManager{
			clm: dm.Get(cloud.DependencyCloudManager).(cloud.Manager),
		}
		smsManager.(*snsManager).initialize()
	default:
		err := errors.New("sms provider unknown")
		logger.Log.Fatal().Err(err).Str("provider", provider)
	}
	dm.Register(DependencySmsManager, smsManager)
}

func ParsePhoneNumber(mobileNumber string) (*phonenumbers.PhoneNumber, error) {
	return phonenumbers.Parse(mobileNumber, "")
}

func GetCharacterCountForMessage(message string) int {
	messageLength := 1
	// We need to identify if the message is unicode, and apply the appropriate character count
	// See https://tinyurl.com/kcvp24d6 for details
	buf := []byte(message)
	// If the number of bytes is not equal to the number of runes, it's a unicode message
	if len(buf) == utf8.RuneCount(buf) {
		// In case the message is greater than 160 characters, we need to split it into multiple messages
		// See https://tinyurl.com/kcvp24d6 for details
		if len(message) > singleSmsCharacterCount {
			msgLength := float64(len(message)) / multiSmsCharacterCount
			messageLength = int(math.Ceil(msgLength))
		}
	} else {
		// In case the message is greater than 70 characters, we need to split it into multiple messages
		// See https://tinyurl.com/kcvp24d6 for details
		if len(message) > unicodeSingleSmsCharacterCount {
			msgLength := float64(len(message)) / unicodeMultiSmsCharacterCount
			messageLength = int(math.Ceil(msgLength))
		}
	}
	return messageLength
}
