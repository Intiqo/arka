package secrets

import (
	"fmt"
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
	cm := dm.Get(config.DependencyConfigManager).(config.Manager)
	cloud.Bootstrap(cloud.ProviderAws)
	secretName := cm.GetValueForKey(AwsSecretNameKey)
	Bootstrap(ProviderAwsSecrets, secretName)
	fmt.Println("Secret Name: ", secretName)
	fmt.Printf("Bootstrap Successful (Aws)!")
	ts.m = dm.Get(DependencySecretsManager).(Manager)
}

func TestAwsSecretsManager(t *testing.T) {
	suite.Run(t, new(AwsSecretsManagerTestSuite))
}

func (ts AwsSecretsManagerTestSuite) Test_awsSecretsManager_GetValueForKey() {
	ts.Run("success", func() {
		res := ts.m.GetValueForKey(constants.AppNameKey)
		fmt.Printf("Get Key result(App Name): %s", res)
		assert.Equal(ts.T(), "App", res)
	})

	ts.Run("unknown key", func() {
		res := ts.m.GetValueForKey("unknown")
		fmt.Printf("Get Key result (unknown): %s", res)
		assert.Equal(ts.T(), "", res)
	})
}
