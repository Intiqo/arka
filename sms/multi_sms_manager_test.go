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
	"github.com/adwitiyaio/arka/secrets"
)

type MultiSmsManagerTestSuite struct {
	suite.Suite

	m Manager
}

func TestMultiSmsManager(t *testing.T) {
	suite.Run(t, new(MultiSmsManagerTestSuite))
}

func (ts *MultiSmsManagerTestSuite) SetupSuite() {
	config.Bootstrap(config.ProviderEnvironment, "../test.env")
	secrets.Bootstrap(secrets.ProviderEnvironment, "")
	cloud.Bootstrap(cloud.ProviderAws)
	err := os.Setenv("CI", "true")
	require.NoError(ts.T(), err)
	Bootstrap()
	ts.m = dependency.GetManager().Get(DependencySmsManager).(Manager)
}

func (ts *MultiSmsManagerTestSuite) Test_multiSmsManager_SendSms() {
	ts.Run("success - click send", func() {
		options := Options{
			Provider:   ProviderMulti,
			Recipients: []string{"+91 9109101910", "+91 9209101920"},
			Message:    "You mustn't be afraid to dream a little bigger darling!",
		}

		res, err := ts.m.SendSms(options)
		assert.NoError(ts.T(), err)
		assert.Nil(ts.T(), res)
	})

	ts.Run("success - smsbroadcast", func() {
		options := Options{
			Provider:   ProviderMulti,
			Recipients: []string{"+61 450 780 453", "+61 418 559 764"},
			Message:    "You mustn't be afraid to dream a little bigger darling!",
		}

		res, err := ts.m.SendSms(options)
		assert.NoError(ts.T(), err)
		assert.Nil(ts.T(), res)
	})

}
