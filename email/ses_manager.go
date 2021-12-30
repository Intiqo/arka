package email

import (
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

func (sm *sesManager) SendEmail(options Options) error {
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

func (sm *sesManager) dispatch(input *ses.SendEmailInput) error {
	if os.Getenv("CI") == "true" {
		return nil
	}
	result, err := sm.ses.SendEmail(input)
	if err != nil {
		if emailError, ok := err.(awserr.Error); ok {
			logger.Log.Error().Err(emailError).Msgf("failed to send email, reason: %s", emailError.Code())
		} else {
			logger.Log.Error().Err(err).Msgf("failed to send email")
		}
		return err
	}

	logger.Log.Debug().Strs("data", []string{input.Message.Subject.GoString(), result.GoString()}).Msg("sent email successfully")
	return nil
}

func (sm *sesManager) initialize() {
	session := sm.clm.GetSession()
	sm.ses = ses.New(session)
	sm.region = sm.clm.GetRegion()
}
