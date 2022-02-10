package email

import (
	"bytes"
	"gopkg.in/gomail.v2"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ses"

	"github.com/adwitiyaio/arka/cloud"
	"github.com/adwitiyaio/arka/config"
	"github.com/adwitiyaio/arka/logger"
)

type sesManager struct {
	cm     config.Manager
	clm    cloud.Manager
	region string
	ses    *ses.SES
}

func (sm sesManager) SendEmail(options Options) error {
	if len(options.Attachments) > 0 {
		return sm.sendRawEmail(options)
	}
	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			BccAddresses: aws.StringSlice(options.Bcc),
			CcAddresses:  aws.StringSlice(options.Cc),
			ToAddresses:  aws.StringSlice(options.To),
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String("UTF-8"),
					Data:    aws.String(options.Html),
				},
				Text: &ses.Content{
					Charset: aws.String("UTF-8"),
					Data:    aws.String(options.Text),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String("UTF-8"),
				Data:    aws.String(options.Subject),
			},
		},
		Source:           aws.String(options.Sender),
		ReplyToAddresses: aws.StringSlice([]string{options.ReplyToAddress}),
	}
	return sm.dispatch(input)
}

func (sm sesManager) dispatch(input *ses.SendEmailInput) error {
	if os.Getenv("CI") == "true" {
		return nil
	}
	_, err := sm.ses.SendEmail(input)
	if err != nil {
		if emailError, ok := err.(awserr.Error); ok {
			logger.Log.Error().Err(emailError).Msgf("failed to send email, reason: %s", emailError.Code())
		} else {
			logger.Log.Error().Err(err).Msgf("failed to send email")
		}
		return err
	}

	logger.Log.Debug().Msg("sent email successfully")
	return nil
}

func (sm sesManager) sendRawEmail(options Options) error {
	// Create raw message
	msg := gomail.NewMessage()

	// Set to
	var recipients []*string
	for _, r := range options.To {
		recipient := r
		recipients = append(recipients, &recipient)
	}

	// Set to emails
	msg.SetHeader("To", options.To...)

	// Set cc
	if len(options.Cc) > 0 {
		for _, r := range options.Cc {
			recipient := r
			recipients = append(recipients, &recipient)
		}
		msg.SetHeader("cc", options.Cc...)
	}

	// Set Bcc
	if len(options.Bcc) > 0 {
		for _, r := range options.Bcc {
			recipient := r
			recipients = append(recipients, &recipient)
		}
		msg.SetHeader("bcc", options.Bcc...)
	}

	msg.SetAddressHeader("From", options.Sender, options.Sender)
	msg.SetHeader("To", options.To...)
	msg.SetHeader("Subject", options.Subject)
	msg.SetBody("text/html", options.Html)

	if len(options.Attachments) > 0 {
		for _, f := range options.Attachments {
			msg.Attach(f)
		}
	}

	// create a new buffer to add raw data
	var emailRaw bytes.Buffer
	_, err := msg.WriteTo(&emailRaw)
	if err != nil {
		return err
	}

	// Create new raw message
	message := ses.RawMessage{Data: emailRaw.Bytes()}
	input := &ses.SendRawEmailInput{Source: &options.Sender, Destinations: recipients, RawMessage: &message}
	return sm.dispatchRawEmail(input)
}

func (sm sesManager) dispatchRawEmail(input *ses.SendRawEmailInput) error {
	if os.Getenv("CI") == "true" {
		return nil
	}
	_, err := sm.ses.SendRawEmail(input)
	if err != nil {
		if emailError, ok := err.(awserr.Error); ok {
			logger.Log.Error().Err(emailError).Msgf("failed to send email, reason: %s", emailError.Code())
		} else {
			logger.Log.Error().Err(err).Msgf("failed to send email")
		}
		return err
	}

	logger.Log.Debug().Msg("sent raw email successfully")
	return nil
}

func (sm *sesManager) initialize() {
	session := sm.clm.GetSession()
	sm.ses = ses.New(session)
	sm.region = sm.clm.GetRegion()
}
