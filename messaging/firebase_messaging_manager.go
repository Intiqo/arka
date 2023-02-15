package messaging

import (
	"context"
	"fmt"
	"os"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"google.golang.org/api/option"

	"github.com/adwitiyaio/arka/logger"
	"github.com/adwitiyaio/arka/secrets"
)

const messagingFirebaseConfigPath = "FIREBASE_MESSAGING_CONFIG_PATH"

type messageResponse struct {
	MessageID string `json:"message_id"`
	Error     error  `json:"error"`
	Success   bool   `json:"success"`
}

type firebaseManager struct {
	sm     secrets.Manager
	client *messaging.Client
}

func (m *firebaseManager) initialize() {
	fbConfigPath := m.sm.GetValueForKey(messagingFirebaseConfigPath)
	opt := option.WithCredentialsFile(fbConfigPath)
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("failed to initialize firebase app")
	}
	ctx := context.Background()
	if os.Getenv("FIREBASE_MESSAGING_DISABLED") == "true" {
		return
	}
	m.client, err = app.Messaging(ctx)
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("failed to get firebase messaging client")
	}
}

func (m *firebaseManager) SendNotification(message Message) (interface{}, []string, error) {
	dataMap := make(map[string]string)
	for key, value := range message.Data {
		strValue := fmt.Sprintf("%v", value)
		dataMap[key] = strValue
	}
	msg := &messaging.MulticastMessage{
		Tokens: message.Tokens,
		Data:   dataMap,
		Notification: &messaging.Notification{
			Title:    message.Title,
			Body:     message.Body,
			ImageURL: message.ImageUrl,
		},
		Android: &messaging.AndroidConfig{
			Data: dataMap,
			Notification: &messaging.AndroidNotification{
				Title:     message.Title,
				Body:      message.Body,
				Sound:     "sound1",
				ChannelID: "shift_channel",
				ImageURL:  "",
			},
		},
		APNS: &messaging.APNSConfig{
			Payload: &messaging.APNSPayload{
				Aps: &messaging.Aps{
					Sound: "default",
				},
			},
		},
	}

	br, err := m.client.SendMulticast(context.Background(), msg)
	if err != nil {
		logger.Log.Error().Err(err).Msgf("failed to send multicast message")
		return nil, nil, err
	}

	if br.SuccessCount > 0 {
		logger.Log.Info().Msgf("successfully sent message to %d devices", br.SuccessCount)
	}
	if br.FailureCount > 0 {
		logger.Log.Info().Msgf("failed to send message to %d devices", br.FailureCount)
	}

	failedTokens := make([]string, 0)
	if br.FailureCount > 0 {
		for idx, resp := range br.Responses {
			if !resp.Success {
				// The order of responses corresponds to the order of the registration tokens.
				failedTokens = append(failedTokens, message.Tokens[idx])
				logger.Log.Error().Err(resp.Error).Msgf("failed to send push notification to device")
			}
		}
	}
	responseData := make([]messageResponse, 0)
	for _, resp := range br.Responses {
		responseData = append(responseData, messageResponse{
			MessageID: resp.MessageID,
			Error:     resp.Error,
			Success:   resp.Success,
		})
	}
	return responseData, failedTokens, nil
}
