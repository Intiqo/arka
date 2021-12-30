package cloud

import (
	"testing"

	"github.com/adwitiyaio/arka/config"
	"github.com/adwitiyaio/arka/dependency"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
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

func (ts AwsManagerTestSuite) Test_awsCloudManager_GetSession() {
	ts.Run("success", func() {
		sess := ts.clm.GetSession()
		assert.NotNil(ts.T(), sess)
	})
}

func (ts AwsManagerTestSuite) Test_awsCloudManager_GetRegion() {
	ts.Run("success", func() {
		region := ts.clm.GetRegion()
		assert.NotNil(ts.T(), region)
	})
}
