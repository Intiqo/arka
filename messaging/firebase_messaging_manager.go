package messaging

import (
	"context"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"google.golang.org/api/option"

	"github.com/adwitiyaio/arka/config"
	"github.com/adwitiyaio/arka/logger"
)

const messagingFirebaseConfigPath = "FIREBASE_MESSAGING_CONFIG_PATH"

type firebaseManager struct {
	cm     config.Manager
	client *messaging.Client
}

func (f *firebaseManager) initialize() {
	fbConfigPath := f.cm.GetValueForKey(messagingFirebaseConfigPath)
	opt := option.WithCredentialsFile(fbConfigPath)
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("failed to initialize firebase app")
	}
	ctx := context.Background()
	f.client, err = app.Messaging(ctx)
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("failed to get firebase messaging client")
	}
}

func (f firebaseManager) SendNotification(message Message) (interface{}, []string) {
	msg := &messaging.MulticastMessage{
		Tokens: message.Tokens,
		Data:   message.Data,
		Notification: &messaging.Notification{
			Title:    message.Title,
			Body:     message.Body,
			ImageURL: message.ImageUrl,
		},
		Android: &messaging.AndroidConfig{
			Data: message.Data,
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

	br, err := f.client.SendMulticast(context.Background(), msg)
	if err != nil {
		logger.Log.Error().Err(err).Msgf("failed to send multicast message")
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
				logger.Log.Error().Err(resp.Error).Msgf("failed to send push notification")
			}
		}
	}
	return br, failedTokens
}
