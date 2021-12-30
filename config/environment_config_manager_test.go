package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/adwitiyaio/arka/constants"
	"github.com/adwitiyaio/arka/dependency"
)

type EnvironmentConfigManagerTestSuite struct {
	suite.Suite
	m Manager
}

func (ts *EnvironmentConfigManagerTestSuite) SetupSuite() {
	Bootstrap(ProviderEnvironment, "../test.env")
	ts.m = dependency.GetManager().Get(DependencyConfigManager).(Manager)
}

func TestEnvironmentConfigManager(t *testing.T) {
	suite.Run(t, new(EnvironmentConfigManagerTestSuite))
}

func (ts EnvironmentConfigManagerTestSuite) Test_environmentConfigManager_GetValueForKey() {
	ts.Run("success", func() {
		res := ts.m.GetValueForKey(constants.AppNameKey)
		assert.Equal(ts.T(), "App", res)
	})
}
