package scheduler

import (
	"github.com/adwitiyaio/arka/dependency"
	"github.com/robfig/cron/v3"
)

const DependencyScheduleManager = "schedule_manager"

const ProviderCron = "cron"

type Manager interface {
	// AddFunc ... Registers a scheduler
	AddFunc(interval string, handler func()) (cron.EntryID, error)
}

// Bootstrap ... Bootstraps the schedule manager
func Bootstrap(provider string) {
	d := dependency.GetManager()
	var c interface{}
	switch provider {
	case ProviderCron:
		c = &cronManager{
			cron: cron.New(),
		}
		c.(*cronManager).cron.Start()
	}
	d.Register(DependencyScheduleManager, c)
}
