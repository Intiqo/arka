package monitoring

import (
	"github.com/adwitiyaio/arka/config"
	"github.com/adwitiyaio/arka/dependency"
)

const DependencyMonitoringManager = "monitoring_manager"

const ProviderNewRelic = "new_relic"

type Manager interface{}

// Bootstrap ... Bootstraps the schedule manager
func Bootstrap(provider string) {
	dm := dependency.GetManager()
	d := dependency.GetManager()
	var nrm interface{}
	switch provider {
	case ProviderNewRelic:
		nrm = &newRelicManager{
			cm: dm.Get(config.DependencyConfigManager).(config.Manager),
		}
		nrm.(*newRelicManager).initialize()
	}
	d.Register(DependencyMonitoringManager, nrm)
}
