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

type MultiMessagingManagerTestSuite struct {
	suite.Suite

	m Manager
}

func TestMultiMessagingManager(t *testing.T) {
	suite.Run(t, new(MultiMessagingManagerTestSuite))
}

func (ts *MultiMessagingManagerTestSuite) SetupSuite() {
	config.Bootstrap(config.ProviderEnvironment, "../test.env")
	secrets.Bootstrap(secrets.ProviderEnvironment, "")
	err := os.Setenv("CI", "true")
	require.NoError(ts.T(), err)
	Bootstrap()
	ts.m = dependency.GetManager().Get(DependencyMessagingManager).(Manager)
}

func (ts *MultiMessagingManagerTestSuite) Test_multiMessagingManager_SendNotificationWithProvider() {
	ts.Run("success - firebase - invalid tokens", func() {
		message := Message{
			Title:    gofakeit.JobTitle(),
			Body:     gofakeit.JobDescriptor(),
			ImageUrl: gofakeit.ImageURL(20, 30),
			Data:     map[string]interface{}{"test": "test"},
			Tokens:   []string{gofakeit.UUID()},
			Channel:  "test",
		}

		res, failedTokens, err := ts.m.SendNotificationWithProvider(message, ProviderFirebase)
		require.NoError(ts.T(), err)
		assert.Equal(ts.T(), message.Tokens, failedTokens)
		assert.False(ts.T(), res.([]messageResponse)[0].Success)
	})

	ts.Run("success - onesignal - invalid tokens", func() {
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
