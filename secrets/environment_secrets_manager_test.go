package secrets

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/adwitiyaio/arka/config"
	"github.com/adwitiyaio/arka/constants"
	"github.com/adwitiyaio/arka/dependency"
)

type EnvironmentSecretsManagerTestSuite struct {
	suite.Suite
	m Manager
}

func (ts *EnvironmentSecretsManagerTestSuite) SetupSuite() {
	dm := dependency.GetManager()
	config.Bootstrap(config.ProviderEnvironment, "../test.env")
	Bootstrap(ProviderEnvironment, "")
	ts.m = dm.Get(DependencySecretsManager).(Manager)
}

func TestEnvironmentSecretsManager(t *testing.T) {
	suite.Run(t, new(EnvironmentSecretsManagerTestSuite))
}

func (ts EnvironmentSecretsManagerTestSuite) Test_environmentSecretsManager_GetValueForKey() {
	ts.Run("success", func() {
		res := ts.m.GetValueForKey(constants.AppNameKey)
		assert.Equal(ts.T(), "App", res)
	})

	ts.Run("unknown key", func() {
		res := ts.m.GetValueForKey("unknown")
		assert.Equal(ts.T(), "", res)
	})
}
