package sms

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/adwitiyaio/arka/cloud"
	"github.com/adwitiyaio/arka/config"
	"github.com/adwitiyaio/arka/dependency"
)

type SnsManagerTestSuite struct {
	suite.Suite

	m Manager
}

func TestSnsManager(t *testing.T) {
	suite.Run(t, new(SnsManagerTestSuite))
}

func (ts *SnsManagerTestSuite) SetupSuite() {
	err := os.Setenv("CI", "true")
	config.Bootstrap(config.ProviderEnvironment, "../test.env")
	cloud.Bootstrap(cloud.ProviderAws)
	require.NoError(ts.T(), err)
	Bootstrap(ProviderSns)
	ts.m = dependency.GetManager().Get(DependencySmsManager).(Manager)
}

func (ts SnsManagerTestSuite) Test_snsManager_SendSms() {
	ts.Run("success", func() {
		options := Options{
			Recipients: []string{"+91 9109101910", "+91 9209101920"},
			Message:    "You mustn't be afraid to dream a little bigger darling!",
		}

		res, err := ts.m.SendSms(options)
		assert.NoError(ts.T(), err)
		assert.Nil(ts.T(), res)
	})
}
