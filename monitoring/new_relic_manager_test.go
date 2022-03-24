package monitoring

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/adwitiyaio/arka/config"
	"github.com/adwitiyaio/arka/dependency"
)

type MonitoringManagerTestSuite struct {
	suite.Suite

	mm Manager
}

func (ts *MonitoringManagerTestSuite) SetupSuite() {
	config.Bootstrap(config.ProviderEnvironment, "../test.env")
	Bootstrap(ProviderNewRelic)
	ts.mm = dependency.GetManager().Get(DependencyMonitoringManager).(Manager)
}

func TestMonitoringManager(t *testing.T) {
	suite.Run(t, new(MonitoringManagerTestSuite))
}

func (ts *MonitoringManagerTestSuite) Test_StartMonitoring() {
	ts.Run("success", func() {
	})
}
