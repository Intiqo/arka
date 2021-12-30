package sms

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sns"

	"github.com/adwitiyaio/arka/cloud"
	"github.com/adwitiyaio/arka/logger"
)

type snsManager struct {
	clm  cloud.Manager
	snss *sns.SNS
}

func (snsm *snsManager) initialize() {
	session := snsm.clm.GetSession()
	snsm.snss = sns.New(session)
}

func (snsm snsManager) SendSms(options Options) error {
	for _, recipient := range options.Recipients {
		input := &sns.PublishInput{
			Message:     aws.String(options.Message),
			PhoneNumber: aws.String(recipient),
		}
		if os.Getenv("CI") == "true" {
			return nil
		}
		_, err := snsm.snss.Publish(input)

		if err != nil {
			logger.Log.Error().Err(err).Msgf("failed to send SMS for the mobile number, %s", recipient)
			continue
		}
	}

	return nil
}
