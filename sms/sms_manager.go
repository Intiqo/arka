package sms

import (
	"errors"
	"fmt"
	"math"
	"unicode/utf8"

	"github.com/nyaruka/phonenumbers"

	"github.com/adwitiyaio/arka/cloud"
	"github.com/adwitiyaio/arka/dependency"
	"github.com/adwitiyaio/arka/logger"
	"github.com/adwitiyaio/arka/secrets"
)

const DependencySmsManager = "sms_manager"

const singleSmsCharacterCount = 160
const multiSmsCharacterCount = 153
const unicodeSingleSmsCharacterCount = 70
const unicodeMultiSmsCharacterCount = 66

const ProviderMulti = "multi"
const ProviderSns = "sns"
const ProviderTermii = "termii"
const ProviderBurstSms = "burstsms"

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
	SendSms(options Options) (interface{}, error)
}

// Bootstrap ... Bootstraps the SMS Manager
func Bootstrap(provider string) {
	dm := dependency.GetManager()
	var smsManager interface{}
	switch provider {
	case ProviderMulti:
		smsManager = &multiSmsManager{
			sm: dm.Get(secrets.DependencySecretsManager).(secrets.Manager),
		}
		smsManager.(*multiSmsManager).initialize()
	case ProviderSns:
		smsManager = &snsManager{
			clm: dm.Get(cloud.DependencyCloudManager).(cloud.Manager),
		}
		smsManager.(*snsManager).initialize()
	case ProviderTermii:
		smsManager = &termiiManager{
			sm: dm.Get(secrets.DependencySecretsManager).(secrets.Manager),
		}
		smsManager.(*termiiManager).initialize()
	case ProviderBurstSms:
		smsManager = &BurstSmsManager{
			sm: dm.Get(secrets.DependencySecretsManager).(secrets.Manager),
		}
		smsManager.(*BurstSmsManager).initialize()
	default:
		err := errors.New("sms provider unknown")
		logger.Log.Fatal().Err(err).Str("provider", provider)
	}
	dm.Register(DependencySmsManager, smsManager)
}

func ParsePhoneNumber(mobileNumber string) (*phonenumbers.PhoneNumber, error) {
	return phonenumbers.Parse(mobileNumber, "")
}

// NormalizePhoneNumber ... Normalizes a phone number by removing leading zeros
// Returns the normalized phone number and the country code
func NormalizePhoneNumber(phoneNumber string) (string, string) {
	var countryCode string
	var normalizedPhoneNumber string
	// We need to normalise the phone number to remove any leading 0's
	p, err := ParsePhoneNumber(phoneNumber)
	if err == nil {
		cc := p.GetCountryCode()
		nationalNumber := p.GetNationalNumber()
		normalizedPhoneNumber = fmt.Sprintf("+%d%d", cc, nationalNumber)
		countryCode = fmt.Sprintf("+%d", cc)
		return normalizedPhoneNumber, countryCode
	}
	return phoneNumber, countryCode
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
