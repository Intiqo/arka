package messaging

import (
	"os"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/adwitiyaio/arka/config"
	"github.com/adwitiyaio/arka/dependency"
	"github.com/adwitiyaio/arka/secrets"
)

type OneSignalManagerTestSuite struct {
	suite.Suite

	m Manager
}

func TestOneSignalManager(t *testing.T) {
	suite.Run(t, new(OneSignalManagerTestSuite))
}

func (ts *OneSignalManagerTestSuite) SetupSuite() {
	config.Bootstrap(config.ProviderEnvironment, "../test.env")
	secrets.Bootstrap(secrets.ProviderEnvironment, "")
	err := os.Setenv("CI", "true")
	require.NoError(ts.T(), err)
	Bootstrap()
	ts.m = dependency.GetManager().Get(DependencyMessagingManager).(Manager)
}

func (ts *OneSignalManagerTestSuite) Test_oneSignalManager_SendNotification() {
	ts.Run("success - invalid tokens", func() {
		message := Message{
			Title:    gofakeit.JobTitle(),
			Body:     gofakeit.JobDescriptor(),
			ImageUrl: gofakeit.ImageURL(20, 30),
			Data:     map[string]interface{}{"test": "test"},
			Tokens:   []string{gofakeit.UUID()},
		}

		res, failedTokens, err := ts.m.SendNotificationWithProvider(message, ProviderOneSignal)
		require.NoError(ts.T(), err)
		assert.Equal(ts.T(), len(message.Tokens), len(failedTokens))
		assert.NotNil(ts.T(), res)
	})
}
