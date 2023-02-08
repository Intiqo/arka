package cloud

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/adwitiyaio/arka/config"
	"github.com/adwitiyaio/arka/dependency"
)

type AwsManagerTestSuite struct {
	suite.Suite

	clm Manager
}

func TestCloudManager(t *testing.T) {
	suite.Run(t, new(AwsManagerTestSuite))
}

func (ts *AwsManagerTestSuite) SetupSuite() {
	dm := dependency.GetManager()
	config.Bootstrap(config.ProviderEnvironment, "../test.env")
	Bootstrap(ProviderAws)
	ts.clm = dm.Get(DependencyCloudManager).(Manager)
}

func (ts AwsManagerTestSuite) Test_awsCloudManager_GetConfig() {
	ts.Run("success", func() {
		config := ts.clm.GetConfig()
		assert.NotNil(ts.T(), config)
	})
}

func (ts AwsManagerTestSuite) Test_awsCloudManager_GetRegion() {
	ts.Run("success", func() {
		region := ts.clm.GetRegion()
		assert.NotNil(ts.T(), region)
	})
}
