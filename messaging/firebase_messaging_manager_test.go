package messaging

import (
	"os"
	"testing"

	"firebase.google.com/go/v4/messaging"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/adwitiyaio/arka/config"
	"github.com/adwitiyaio/arka/dependency"
)

type FirebaseMessagingManagerTestSuite struct {
	suite.Suite

	m Manager
}

func TestFirebaseMessagingManager(t *testing.T) {
	suite.Run(t, new(FirebaseMessagingManagerTestSuite))
}

func (ts *FirebaseMessagingManagerTestSuite) SetupSuite() {
	config.Bootstrap(config.ProviderEnvironment, "../test.env")
	err := os.Setenv("CI", "true")
	require.NoError(ts.T(), err)
	Bootstrap(ProviderFirebase)
	ts.m = dependency.GetManager().Get(DependencyMessagingManager).(Manager)
}

func (ts FirebaseMessagingManagerTestSuite) Test_firebaseMessagingManager_SendNotification() {
	ts.Run("success - invalid tokens", func() {
		message := Message{
			Title:    gofakeit.JobTitle(),
			Body:     gofakeit.JobDescriptor(),
			ImageUrl: gofakeit.ImageURL(20, 30),
			Data:     map[string]string{"test": "test"},
			Tokens:   []string{gofakeit.UUID()},
		}

		res, failedTokens := ts.m.SendNotification(message)
		assert.Equal(ts.T(), message.Tokens, failedTokens)
		assert.Equal(ts.T(), 1, res.(*messaging.BatchResponse).FailureCount)
	})
}
