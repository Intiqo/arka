package sms

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"

	"github.com/adwitiyaio/arka/cloud"
	"github.com/adwitiyaio/arka/logger"
)

type snsManager struct {
	clm    cloud.Manager
	client *sns.Client
}

func (snsm *snsManager) initialize() {
	config := snsm.clm.GetConfig()
	snsm.client = sns.NewFromConfig(config)
}

func (snsm snsManager) SendSms(options Options) (interface{}, error) {
	var res []*sns.PublishOutput
	for _, recipient := range options.Recipients {
		input := &sns.PublishInput{
			Message:     aws.String(options.Message),
			PhoneNumber: aws.String(recipient),
		}
		if os.Getenv("CI") == "true" {
			return nil, nil
		}
		out, err := snsm.client.Publish(context.TODO(), input)
		if err != nil {
			logger.Log.Error().Err(err).Msgf("failed to send SMS for the mobile number, %s", recipient)
			continue
		}
		res = append(res, out)
	}

	return res, nil
}
