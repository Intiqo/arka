package sms

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/adwitiyaio/arka/config"
	"github.com/adwitiyaio/arka/dependency"
)

type TermiiManagerTestSuite struct {
	suite.Suite

	m Manager
}

func TestTermiiManager(t *testing.T) {
	suite.Run(t, new(TermiiManagerTestSuite))
}

func (ts *TermiiManagerTestSuite) SetupSuite() {
	err := os.Setenv("CI", "true")
	config.Bootstrap(config.ProviderEnvironment, "../test.env")
	require.NoError(ts.T(), err)
	Bootstrap(ProviderTermii)
	ts.m = dependency.GetManager().Get(DependencySmsManager).(Manager)
}

func (ts TermiiManagerTestSuite) Test_termiiManager_SendSms() {
	ts.Run(
		"success", func() {
			options := Options{
				Recipients: []string{"+91 9109101910", "+91 9209101920"},
				Message:    "You mustn't be afraid to dream a little bigger darling!",
			}

			res, err := ts.m.SendSms(options)
			assert.NoError(ts.T(), err)
			assert.NotNil(ts.T(), res)
		},
	)
}
