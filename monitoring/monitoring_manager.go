package monitoring

import (
	"github.com/newrelic/go-agent/v3/newrelic"

	"github.com/adwitiyaio/arka/config"
	"github.com/adwitiyaio/arka/dependency"
)

const DependencyMonitoringManager = "monitoring_manager"

const ProviderNewRelic = "new_relic"

type Manager interface {
	StartMonitoring() (*newrelic.Application, error)
}

// Bootstrap ... Bootstraps the schedule manager
func Bootstrap(provider string) {
	dm := dependency.GetManager()
	d := dependency.GetManager()
	var c interface{}
	switch provider {
	case ProviderNewRelic:
		c = &newRelicManager{
			cm: dm.Get(config.DependencyConfigManager).(config.Manager),
		}
	}
	d.Register(DependencyMonitoringManager, c)
}
