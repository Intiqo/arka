package email

import (
	"context"
	"os"
	"time"

	"github.com/mailgun/mailgun-go/v4"

	"github.com/adwitiyaio/arka/config"
	"github.com/adwitiyaio/arka/logger"
)

const domainKey = "MAILGUN_DOMAIN"
const apiKey = "MAILGUN_API_KEY"

type mailgunManager struct {
	cm config.Manager
	mg *mailgun.MailgunImpl
}

func (m *mailgunManager) initialize() {
	domain := m.cm.GetValueForKey(domainKey)
	apiKey := m.cm.GetValueForKey(apiKey)

	if domain == "" || apiKey == "" {
		logger.Log.Panic().Msg("failed to initialize mailgun")
	}

	m.mg = mailgun.NewMailgun(domain, apiKey)
}

func (m *mailgunManager) SendEmail(options Options) (interface{}, error) {
	message := m.mg.NewMessage(options.Sender, options.Subject, options.Text)
	for _, to := range options.To {
		err := message.AddRecipient(to)
		if err != nil {
			logger.Log.Error().Str("recipient", to).Err(err).Stack().Msg("failed to add recipient")
		}
	}

	for _, cc := range options.Cc {
		message.AddCC(cc)
	}

	for _, bcc := range options.Bcc {
		message.AddBCC(bcc)
	}

	message.SetHtml(options.Html)
	if len(options.Attachments) > 0 {
		for _, attachment := range options.Attachments {
			message.AddAttachment(attachment)
		}
	}
	message.SetTracking(true)
	message.SetTrackingClicks(true)
	message.SetTrackingOpens(true)

	return m.dispatch(message, options)
}

func (m *mailgunManager) dispatch(message *mailgun.Message, options Options) (interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	var resp string
	var err error
	if os.Getenv("CI") != "true" {
		resp, _, err = m.mg.Send(ctx, message)
	}

	if err != nil {
		logger.Log.Error().Err(err).Stack().Msg("failed to send email")
		return nil, err
	}
	return resp, nil
}
