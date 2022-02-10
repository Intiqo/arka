package email

import (
	"errors"

	"github.com/adwitiyaio/arka/cloud"
	"github.com/adwitiyaio/arka/config"
	"github.com/adwitiyaio/arka/dependency"
	"github.com/adwitiyaio/arka/logger"
)

const DependencyEmailManager = "email_manager"

const ProviderMailgun = "MAILGUN"
const ProviderSes = "SES"

// Options ... Defines the structure for the fields of an email
type Options struct {
	Sender         string
	Subject        string
	Html           string
	Text           string
	To             []string
	Cc             []string
	Bcc            []string
	Attachments    []string
	ReplyToAddress string
}

// Manager ... An email service to send or receive emails
type Manager interface {
	// SendEmail ... Send email with options
	SendEmail(options Options) error
}

// Bootstrap ... Bootstrap the email service
func Bootstrap(provider string) {
	dm := dependency.GetManager()
	var m interface{}
	switch provider {
	case ProviderMailgun:
		m = &mailgunManager{
			cm: dm.Get(config.DependencyConfigManager).(config.Manager),
		}
		m.(*mailgunManager).initialize()
	case ProviderSes:
		m = &sesManager{
			cm:  dm.Get(config.DependencyConfigManager).(config.Manager),
			clm: dm.Get(cloud.DependencyCloudManager).(cloud.Manager),
		}
		m.(*sesManager).initialize()
	default:
		err := errors.New("email provider unknown")
		logger.Log.Fatal().Err(err).Str("provider", provider)
	}
	dm.Register(DependencyEmailManager, m)
}
