package secrets

import (
	"fmt"
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
	fmt.Println("Bootstrap Successful (Env)!")
	ts.m = dm.Get(DependencySecretsManager).(Manager)
}

func TestEnvironmentSecretsManager(t *testing.T) {
	suite.Run(t, new(EnvironmentSecretsManagerTestSuite))
}

func (ts EnvironmentSecretsManagerTestSuite) Test_environmentSecretsManager_GetValueForKey() {
	ts.Run("success", func() {
		res := ts.m.GetValueForKey(constants.AppNameKey)
		fmt.Println("Get Key result(App Name): ", res)
		assert.Equal(ts.T(), "App", res)
	})

	ts.Run("unknown key", func() {
		res := ts.m.GetValueForKey("unknown")
		fmt.Println("Get Key result (unknown): ", res)
		assert.Equal(ts.T(), "", res)
	})
}
