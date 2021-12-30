package scheduler

import (
	"testing"

	"github.com/adwitiyaio/arka/dependency"
	"github.com/robfig/cron/v3"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type ScheduleManagerTestSuite struct {
	suite.Suite

	cc Manager
}

func (ts *ScheduleManagerTestSuite) SetupSuite() {
	Bootstrap(ProviderCron)
	ts.cc = dependency.GetManager().Get(DependencyScheduleManager).(Manager)
}

func TestScheduleManager(t *testing.T) {
	suite.Run(t, new(ScheduleManagerTestSuite))
}

func (ts *ScheduleManagerTestSuite) Test_AddFunc() {
	ts.Run("success", func() {
		id, err := ts.cc.AddFunc("10 * * * *", func() {
			return
		})
		require.NoError(ts.T(), err)
		require.Equal(ts.T(), cron.EntryID(1), id)
	})

	ts.Run("exception for invalid interval", func() {
		id, err := ts.cc.AddFunc("10 * * * * *", func() {
			return
		})
		require.Error(ts.T(), err)
		require.Equal(ts.T(), cron.EntryID(0), id)
	})

	ts.Run("exception for invalid character in interval", func() {
		id, err := ts.cc.AddFunc("10 * * * A", func() {
			return
		})
		require.Error(ts.T(), err)
		require.Equal(ts.T(), cron.EntryID(0), id)
	})
}
