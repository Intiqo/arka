package event

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsCfg "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sns"

	"github.com/adwitiyaio/arka/cloud"
	"github.com/adwitiyaio/arka/logger"
)

type snsEventManager struct {
	topicArnMap map[string]string
	clm         cloud.Manager
	client      *sns.Client
}

func (s *snsEventManager) initialize() {
	s.topicArnMap = make(map[string]string)
	config, err := awsCfg.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic(err)
	}
	s.client = sns.NewFromConfig(config)
}

// Register a topic. This will create a topic if it doesn't exist
func (s *snsEventManager) Register(event string, handler Handler) error {
	out, err := s.client.CreateTopic(context.TODO(), &sns.CreateTopicInput{
		Name: aws.String(event),
	})
	if err != nil {
		return err
	}
	s.topicArnMap[event] = *out.TopicArn
	return nil
}

func (s *snsEventManager) Publish(event string, data interface{}) error {
	payload, err := json.Marshal(data)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Error marshalling data")
		return err
	}
	result, err := s.client.Publish(context.TODO(), &sns.PublishInput{
		Message:  aws.String(string(payload)),
		TopicArn: aws.String(s.topicArnMap[event]),
	})

	if err != nil {
		logger.Log.Error().Err(err).Msg("Error publishing event")
		return err
	}

	logger.Log.Info().Str("event", event).Str("messageId", *result.MessageId).Msg("Event published to SNS")
	return nil
}
