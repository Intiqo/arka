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

type BurstSmsManagerTestSuite struct {
	suite.Suite

	m Manager
}

func TestBurstSmsManager(t *testing.T) {
	suite.Run(t, new(BurstSmsManagerTestSuite))
}

func (bm *BurstSmsManagerTestSuite) SetupSuite() {
	err := os.Setenv("CI", "true")
	config.Bootstrap(config.ProviderEnvironment, "../test.env")
	secrets.Bootstrap(secrets.ProviderEnvironment, "")
	cloud.Bootstrap(cloud.ProviderAws)
	require.NoError(bm.T(), err)
	Bootstrap()
	bm.m = dependency.GetManager().Get(DependencySmsManager).(Manager)
}

func (bm *BurstSmsManagerTestSuite) Test_burstSmsManager_SendSms() {
	bm.Run(
		"success", func() {
			options := Options{
				Provider:   ProviderBurstSms,
				Recipients: []string{"+91 9109101910", "+91 9209101920"},
				Message:    "You mustn't be afraid to dream a little bigger darling!",
			}

			res, err := bm.m.SendSms(options)
			assert.NoError(bm.T(), err)
			assert.Nil(bm.T(), res)
		},
	)
}
