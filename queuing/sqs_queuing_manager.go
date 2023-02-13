package queuing

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"

	"github.com/adwitiyaio/arka/cloud"
	"github.com/adwitiyaio/arka/logger"
	"github.com/adwitiyaio/arka/secrets"
)

const regionKey = "AWS_REGION"

type sqsManager struct {
	sm     secrets.Manager
	clm    cloud.Manager
	client *sqs.Client
	region string
}

func (s *sqsManager) initialize() {
	config := s.clm.GetConfig()
	s.region = s.sm.GetValueForKey(regionKey)
	s.client = sqs.NewFromConfig(config)
}

func (s sqsManager) SendMessage(options SendOptions) error {
	queueUrl, err := s.client.GetQueueUrl(context.TODO(), &sqs.GetQueueUrlInput{
		QueueName: aws.String(options.QueueName),
	})
	if err != nil {
		return err
	}

	body, err := json.Marshal(options.Data)
	if err != nil {
		return err
	}

	_, err = s.client.SendMessage(context.TODO(), &sqs.SendMessageInput{
		MessageBody:    aws.String(string(body)),
		QueueUrl:       queueUrl.QueueUrl,
		MessageGroupId: aws.String(options.GroupId),
	})

	if err != nil {
		return err
	}

	return nil
}

func (s sqsManager) ReceiveMessage(options ReceiveOptions) (ReceiveResponse, error) {
	var response ReceiveResponse
	queueUrl, err := s.client.GetQueueUrl(context.TODO(), &sqs.GetQueueUrlInput{
		QueueName: aws.String(options.QueueName),
	})
	if err != nil {
		return response, err
	}

	if options.WaitTimeSeconds < 1 {
		options.WaitTimeSeconds = 1
	}
	if options.WaitTimeSeconds > 20 {
		options.WaitTimeSeconds = 20
	}

	if options.DelayTimeout == 0 {
		options.DelayTimeout = 20
	}
	if options.DelayTimeout < 5 {
		options.DelayTimeout = 5
	}

	if options.NumberOfMessages < 1 {
		options.NumberOfMessages = 1
	}
	if options.NumberOfMessages > 10 {
		options.NumberOfMessages = 10
	}

	result, err := s.client.ReceiveMessage(context.TODO(), &sqs.ReceiveMessageInput{
		QueueUrl:            queueUrl.QueueUrl,
		MaxNumberOfMessages: int32(options.NumberOfMessages),
		VisibilityTimeout:   int32(options.DelayTimeout),
		WaitTimeSeconds:     int32(options.WaitTimeSeconds),
	})
	if err != nil {
		return response, err
	}

	if len(result.Messages) == 0 {
		return response, nil
	}

	for _, message := range result.Messages {
		var data interface{}
		err = json.Unmarshal([]byte(*message.Body), &data)
		if err != nil {
			logger.Log.Error().Err(err).Str("message_id", *message.MessageId).Msg("Failed to unmarshal message")
			continue
		}
		response.Messages = append(response.Messages, MessageResponse{
			MessageId: *message.MessageId,
			Data:      data,
			Receipt:   *message.ReceiptHandle,
		})
	}

	return response, nil
}

func (s sqsManager) DeleteMessage(options DeleteOptions) error {
	queueUrl, err := s.client.GetQueueUrl(context.TODO(), &sqs.GetQueueUrlInput{
		QueueName: aws.String(options.QueueName),
	})
	if err != nil {
		return err
	}

	_, err = s.client.DeleteMessage(context.TODO(), &sqs.DeleteMessageInput{
		QueueUrl:      queueUrl.QueueUrl,
		ReceiptHandle: aws.String(options.MessageHandle),
	})
	if err != nil {
		return err
	}
	return nil
}
