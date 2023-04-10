package event

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/adwitiyaio/arka/config"
	"github.com/adwitiyaio/arka/dependency"
	"github.com/adwitiyaio/arka/secrets"
)

var message = "Hello"

func greetUser(_ string, data interface{}) {
	message = data.(string)
}

type LocalEventManagerTestSuite struct {
	suite.Suite

	evm Manager
}

func (ts *LocalEventManagerTestSuite) SetupSuite() {
	config.Bootstrap(config.ProviderEnvironment, "../test.env")
	secrets.Bootstrap(secrets.ProviderEnvironment, "")
	err := os.Setenv("CI", "true")
	require.NoError(ts.T(), err)
	Bootstrap(ProviderLocal)
	ts.evm = dependency.GetManager().Get(DependencyEventManager).(Manager)
}

func TestLocalEventManager(t *testing.T) {
	suite.Run(t, new(LocalEventManagerTestSuite))
}

func (ts LocalEventManagerTestSuite) Test_localEventManager_RegisterAndGet() {
	ts.Run("success", func() {
		const eventName = "greet"
		err := ts.evm.Register(eventName, greetUser)
		require.NoError(ts.T(), err)
		assert.Equal(ts.T(), "Hello", message)
		const eventData = "World"
		err = ts.evm.Publish(eventName, eventData)
		require.NoError(ts.T(), err)
		time.Sleep(time.Millisecond * 100)
		assert.Equal(ts.T(), eventData, message)
	})

}
