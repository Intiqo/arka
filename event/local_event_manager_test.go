package event

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
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
	ts.evm = GetManager()
}

func TestLocalEventManager(t *testing.T) {
	suite.Run(t, new(LocalEventManagerTestSuite))
}

func (ts LocalEventManagerTestSuite) Test_localEventManager_RegisterAndGet() {
	ts.Run("success", func() {
		const eventName = "greet"
		ts.evm.Register(eventName, greetUser)
		assert.Equal(ts.T(), "Hello", message)
		const eventData = "World"
		ts.evm.Publish(eventName, eventData)
		time.Sleep(time.Millisecond * 100)
		assert.Equal(ts.T(), eventData, message)
	})

}
