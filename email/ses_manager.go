package email

import (
	"bytes"
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"gopkg.in/gomail.v2"

	"github.com/adwitiyaio/arka/cloud"
	"github.com/adwitiyaio/arka/config"
	"github.com/adwitiyaio/arka/logger"
)

type sesManager struct {
	cm     config.Manager
	clm    cloud.Manager
	region string
	client *ses.Client
}

func (sm sesManager) SendEmail(options Options) (interface{}, error) {
	if len(options.Attachments) > 0 {
		return sm.sendRawEmail(options)
	}
	input := &ses.SendEmailInput{
		Destination: &types.Destination{
			BccAddresses: options.Bcc,
			CcAddresses:  options.Cc,
			ToAddresses:  options.To,
		},
		Message: &types.Message{
			Body: &types.Body{
				Html: &types.Content{
					Charset: aws.String("UTF-8"),
					Data:    aws.String(options.Html),
				},
				Text: &types.Content{
					Charset: aws.String("UTF-8"),
					Data:    aws.String(options.Text),
				},
			},
			Subject: &types.Content{
				Charset: aws.String("UTF-8"),
				Data:    aws.String(options.Subject),
			},
		},
		Source:           aws.String(options.Sender),
		ReplyToAddresses: []string{options.ReplyToAddress},
	}
	return sm.dispatch(input)
}

func (sm sesManager) dispatch(input *ses.SendEmailInput) (interface{}, error) {
	if os.Getenv("CI") == "true" {
		return nil, nil
	}
	resp, err := sm.client.SendEmail(context.TODO(), input)
	if err != nil {
		if emailError, ok := err.(awserr.Error); ok {
			logger.Log.Error().Err(emailError).Msgf("failed to send email, reason: %s", emailError.Code())
		} else {
			logger.Log.Error().Err(err).Msgf("failed to send email")
		}
		return nil, err
	}
	return resp, nil
}

func (sm sesManager) sendRawEmail(options Options) (interface{}, error) {
	// Create raw message
	msg := gomail.NewMessage()

	// Set to
	var recipients []string
	for _, r := range options.To {
		recipients = append(recipients, r)
	}

	// Set to emails
	msg.SetHeader("To", options.To...)

	// Set cc
	if len(options.Cc) > 0 {
		for _, r := range options.Cc {
			recipients = append(recipients, r)
		}
		msg.SetHeader("cc", options.Cc...)
	}

	// Set Bcc
	if len(options.Bcc) > 0 {
		for _, r := range options.Bcc {
			recipients = append(recipients, r)
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
		return nil, err
	}

	// Create new raw message
	message := &types.RawMessage{Data: emailRaw.Bytes()}
	input := &ses.SendRawEmailInput{Source: &options.Sender, Destinations: recipients, RawMessage: message}
	return sm.dispatchRawEmail(input)
}

func (sm sesManager) dispatchRawEmail(input *ses.SendRawEmailInput) (interface{}, error) {
	if os.Getenv("CI") == "true" {
		return nil, nil
	}
	resp, err := sm.client.SendRawEmail(context.TODO(), input)
	if err != nil {
		if emailError, ok := err.(awserr.Error); ok {
			logger.Log.Error().Err(emailError).Msgf("failed to send email, reason: %s", emailError.Code())
		} else {
			logger.Log.Error().Err(err).Msgf("failed to send email")
		}
		return nil, err
	}
	return resp, nil
}

func (sm *sesManager) initialize() {
	config := sm.clm.GetConfig()
	sm.client = ses.NewFromConfig(config)
	sm.region = sm.clm.GetRegion()
}
