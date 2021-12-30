package sms

import (
	"github.com/adwitiyaio/arka/config"
	"github.com/adwitiyaio/arka/dependency"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"os"
	"testing"
)

type SmsManagerTestSuite struct {
	suite.Suite

	m Manager
}

func TestSmsManager(t *testing.T) {
	suite.Run(t, new(SmsManagerTestSuite))
}

func (ts *SmsManagerTestSuite) SetupSuite() {
	err := os.Setenv("CI", "true")
	require.NoError(ts.T(), err)
	config.Bootstrap(config.ProviderEnvironment, "../test.env")
	Bootstrap(ProviderMulti)
	ts.m = dependency.GetManager().Get(DependencySmsManager).(Manager)
}

func (ts SmsManagerTestSuite) Test_GetCharacterCountForMessage() {
	ts.Run("success - non-unicode - single sms", func() {
		const message = "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam"
		count := GetCharacterCountForMessage(message)
		assert.Equal(ts.T(), 1, count)
	})

	ts.Run("success - non-unicode - multiple sms", func() {
		const message = "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum."
		count := GetCharacterCountForMessage(message)
		assert.Equal(ts.T(), 3, count)
	})

	ts.Run("success - unicode - single sms", func() {
		const message = "Lorem ipsum dolor sit amet, consectetur adipiscing elit 💪🏻"
		count := GetCharacterCountForMessage(message)
		assert.Equal(ts.T(), 1, count)
	})

	ts.Run("success - unicode - multiple sms", func() {
		const message = "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt 💪🏻"
		count := GetCharacterCountForMessage(message)
		assert.Equal(ts.T(), 2, count)
	})
}
