package secrets

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/adwitiyaio/arka/cloud"
	"github.com/adwitiyaio/arka/config"
	"github.com/adwitiyaio/arka/constants"
	"github.com/adwitiyaio/arka/dependency"
)

type AwsSecretsManagerTestSuite struct {
	suite.Suite
	m Manager
}

func (ts *AwsSecretsManagerTestSuite) SetupSuite() {
	dm := dependency.GetManager()
	config.Bootstrap(config.ProviderEnvironment, "../test.env")
	cloud.Bootstrap(cloud.ProviderAws)
	Bootstrap(ProviderAwsSecrets, "ci/env")
	ts.m = dm.Get(DependencySecretsManager).(Manager)
}

func TestAwsSecretsManager(t *testing.T) {
	suite.Run(t, new(AwsSecretsManagerTestSuite))
}

func (ts AwsSecretsManagerTestSuite) Test_awsSecretsManager_GetValueForKey() {
	ts.Run("success", func() {
		res := ts.m.GetValueForKey(constants.AppNameKey)
		assert.Equal(ts.T(), "App", res)
	})

	ts.Run("unknown key", func() {
		res := ts.m.GetValueForKey("unknown")
		assert.Equal(ts.T(), "", res)
	})
}
