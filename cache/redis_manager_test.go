package cache

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/adwitiyaio/arka/config"
	"github.com/adwitiyaio/arka/constants"
	"github.com/adwitiyaio/arka/dependency"
	"github.com/adwitiyaio/arka/secrets"
)

type RedisManagerTestSuite struct {
	suite.Suite
	ccm Manager
}

func TestRedisManager(t *testing.T) {
	suite.Run(t, new(RedisManagerTestSuite))
}

func (ts *RedisManagerTestSuite) SetupSuite() {
	dm := dependency.GetManager()
	config.Bootstrap(config.ProviderEnvironment, "../test.env")
	secrets.Bootstrap(secrets.ProviderEnvironment, "")

	// For coverage, set the db to an invalid integer
	err := os.Setenv("REDIS_DATABASE", "invalid")
	require.NoError(ts.T(), err)
	err = os.Setenv("CI", "false")
	require.NoError(ts.T(), err)
	Bootstrap(ProviderRedis)

	// Set correct values now
	config.Bootstrap(config.ProviderEnvironment, "../test.env")
	Bootstrap(ProviderRedis)
	ts.ccm = dm.Get(DependencyCacheManager).(Manager)
}

func (ts RedisManagerTestSuite) Test_redisCacheManager_GetStatus() {
	ts.Run("success", func() {
		status := ts.ccm.GetStatus()
		assert.Equal(ts.T(), constants.SystemStatusUp, status)
	})
}

func (ts RedisManagerTestSuite) Test_redisCacheManager_Set_Get() {
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
