package messaging

import (
	"context"

	"github.com/OneSignal/onesignal-go-api"

	"github.com/adwitiyaio/arka/logger"
	"github.com/adwitiyaio/arka/secrets"
)

const oneSignalAppIdSecretKey = "ONESIGNAL_APP_ID"
const oneSignalAppKeySecretKey = "ONESIGNAL_APP_KEY"

type oneSignalManager struct {
	sm secrets.Manager

	appId  string
	appKey string

	apiClient *onesignal.APIClient
}

func (m *oneSignalManager) initialize() {
	m.appId = m.sm.GetValueForKey(oneSignalAppIdSecretKey)
	m.appKey = m.sm.GetValueForKey(oneSignalAppKeySecretKey)

	configuration := onesignal.NewConfiguration()
	m.apiClient = onesignal.NewAPIClient(configuration)
}

func (m *oneSignalManager) SendNotification(message Message) (interface{}, []string, error) {
	appAuth := context.WithValue(context.Background(), onesignal.AppAuth, m.appKey)

	notification := *onesignal.NewNotification(m.appId)

	notification.SetIncludePlayerIds(message.Tokens)

	notification.SetName("Transactional")

	title := onesignal.NewStringMap()
	title.SetEn(message.Title)
	notification.SetHeadings(*title)

	content := onesignal.NewStringMap()
	content.SetEn(message.Body)
	notification.SetContents(*content)

	// Set custom Android properties
	notification.SetAndroidSound("sound1")

	// Set custom iOS properties
	notification.SetIosSound("default")
	notification.SetIosBadgeType("SetTo")
	notification.SetIosBadgeCount(1)

	// Set the data object
	notification.SetData(message.Data)

	resp, r, err := m.apiClient.DefaultApi.CreateNotification(appAuth).Notification(notification).Execute()

	if err != nil {
		logger.Log.Error().Err(err).Msg("failed to send notification")
		logger.Log.Error().Interface("response", r).Msg("http raw response")
		return r, nil, err
	}

	if resp != nil && resp.Errors != nil && resp.Errors.InvalidIdentifierError != nil && resp.Errors.InvalidIdentifierError.InvalidPlayerIds != nil {
		return resp, resp.Errors.InvalidIdentifierError.InvalidExternalUserIds, nil
	}

	if resp != nil && resp.Errors != nil && resp.Errors.ArrayOfString != nil {
		return resp, *resp.Errors.ArrayOfString, nil
	}

	return resp, []string{}, nil
}
