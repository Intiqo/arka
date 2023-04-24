package cache

import (
	"github.com/adwitiyaio/arka/config"
	"github.com/adwitiyaio/arka/constants"
	"github.com/adwitiyaio/arka/dependency"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"testing"
)

type LocalManagerTestSuite struct {
	suite.Suite
	ccm Manager
}

func TestLocalManager(t *testing.T) {
	suite.Run(t, new(LocalManagerTestSuite))
}

func (ts *LocalManagerTestSuite) SetupSuite() {
	dm := dependency.GetManager()
	config.Bootstrap(config.ProviderEnvironment, "../test.env")
	Bootstrap(ProviderLocal)
	ts.ccm = dm.Get(DependencyCacheManager).(Manager)
}

func (ts *LocalManagerTestSuite) Test_localCacheManager_GetStatus() {
	ts.Run("success", func() {
		status := ts.ccm.GetStatus()
		assert.Equal(ts.T(), constants.SystemStatusUp, status)
	})
}

func (ts *LocalManagerTestSuite) Test_localCacheManager_Set_Get() {
	ts.Run("success", func() {
		key := "test_key"
		val := "Test Value"
		err := ts.ccm.Set(key, val)
		require.NoError(ts.T(), err)
		result, err := ts.ccm.Get(key)
		require.NoError(ts.T(), err)
		assert.Equal(ts.T(), val, result)
	})
}
