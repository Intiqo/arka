package event

import (
	"context"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/adwitiyaio/arka/cloud"
	"github.com/adwitiyaio/arka/config"
	"github.com/adwitiyaio/arka/dependency"
)

type SnsEventManagerTestSuite struct {
	suite.Suite

	evm Manager
}

func (ts *SnsEventManagerTestSuite) SetupSuite() {
	config.Bootstrap(config.ProviderEnvironment, "../test.env")
	cloud.Bootstrap(cloud.ProviderAws)
	err := os.Setenv("CI", "true")
	require.NoError(ts.T(), err)
	Bootstrap(ProviderSns)
	ts.evm = dependency.GetManager().Get(DependencyEventManager).(Manager)
}

func TestSnsEventManager(t *testing.T) {
	suite.Run(t, new(SnsEventManagerTestSuite))
}

func (ts *SnsEventManagerTestSuite) Test_snsEventManager_RegisterAndPublish() {
	ts.Run("success", func() {
		const eventName = "greet"
		err := ts.evm.Register(eventName, nil)
		require.NoError(ts.T(), err)
		eventData := struct {
			Name string `json:"name"`
		}{
			Name: "John",
		}
		err = ts.evm.Publish(eventName, eventData)
		require.NoError(ts.T(), err)
		cl := ts.evm.(*snsEventManager)
		_, err = cl.client.DeleteTopic(context.Background(), &sns.DeleteTopicInput{
			TopicArn: aws.String(cl.topicArnMap[eventName]),
		})
		require.NoError(ts.T(), err)
	})

}
