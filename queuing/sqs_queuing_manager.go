package queuing

import (
	"encoding/json"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"

	"github.com/adwitiyaio/arka/cloud"
	"github.com/adwitiyaio/arka/config"
	"github.com/adwitiyaio/arka/logger"
)

const regionKey = "AWS_REGION"

type sqsManager struct {
	cm      config.Manager
	clm     cloud.Manager
	session *session.Session
	region  string
}

func (s *sqsManager) initialize() {
	s.session = s.clm.GetSession()
	s.region = s.cm.GetValueForKey(regionKey)
}

func (s sqsManager) SendMessage(options SendOptions) error {
	svc := sqs.New(s.session)
	queueUrl, err := svc.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: aws.String(options.QueueName),
	})
	if err != nil {
		return err
	}

	body, err := json.Marshal(options.Data)
	if err != nil {
		return err
	}

	_, err = svc.SendMessage(&sqs.SendMessageInput{
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
	svc := sqs.New(s.session)
	queueUrl, err := svc.GetQueueUrl(&sqs.GetQueueUrlInput{
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

	result, err := svc.ReceiveMessage(&sqs.ReceiveMessageInput{
		QueueUrl:            queueUrl.QueueUrl,
		MaxNumberOfMessages: aws.Int64(options.NumberOfMessages),
		VisibilityTimeout:   aws.Int64(options.DelayTimeout),
		WaitTimeSeconds:     aws.Int64(options.WaitTimeSeconds),
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
	svc := sqs.New(s.session)
	queueUrl, err := svc.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: aws.String(options.QueueName),
	})
	if err != nil {
		return err
	}

	_, err = svc.DeleteMessage(&sqs.DeleteMessageInput{
		QueueUrl:      queueUrl.QueueUrl,
		ReceiptHandle: aws.String(options.MessageHandle),
	})
	if err != nil {
		return err
	}
	return nil
}
