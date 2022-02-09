package queuing

import (
	"testing"

	"github.com/mitchellh/mapstructure"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/adwitiyaio/arka/cloud"
	"github.com/adwitiyaio/arka/config"
	"github.com/adwitiyaio/arka/dependency"
	"github.com/adwitiyaio/arka/logger"
)

type testData struct {
	Message string
}

const queueName = "arka-test.fifo"

type SqsManagerTestSuite struct {
	suite.Suite

	m Manager
}

func TestSqsManager(t *testing.T) {
	suite.Run(t, new(SqsManagerTestSuite))
}

func (ts *SqsManagerTestSuite) SetupSuite() {
	dm := dependency.GetManager()
	config.Bootstrap(config.ProviderEnvironment, "../test.env")
	cloud.Bootstrap(cloud.ProviderAws)
	Bootstrap(ProviderSQS)
	ts.m = dm.Get(DependencyQueuingManager).(Manager)
}

func (ts *SqsManagerTestSuite) Test_sqsManager_SendReceiveAndDeleteMessage() {
	ts.Run("success", func() {
		data := testData{
			Message: "test message",
		}
		err := ts.m.SendMessage(SendOptions{
			Options: Options{QueueName: queueName},
			GroupId: "group-1",
			Data:    data,
		})
		require.NoError(ts.T(), err)

		result, err := ts.m.ReceiveMessage(ReceiveOptions{
			Options:          Options{QueueName: queueName},
			DelayTimeout:     30,
			NumberOfMessages: 1,
		})
		require.NoError(ts.T(), err)
		for _, msg := range result.Messages {
			var messageData testData
			dataMap := msg.Data.(map[string]interface{})
			err = mapstructure.Decode(dataMap, &messageData)
			require.NoError(ts.T(), err)
			assert.NotEmpty(ts.T(), msg.MessageId)
			assert.NotEmpty(ts.T(), msg.Receipt)
			assert.Equal(ts.T(), data.Message, messageData.Message)
			logger.Log.Info().Msgf("data: %v", messageData)
			err = ts.m.DeleteMessage(DeleteOptions{
				Options:       Options{QueueName: queueName},
				MessageHandle: msg.Receipt,
			})
			require.NoError(ts.T(), err)
		}
	})
}
