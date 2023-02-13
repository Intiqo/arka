package monitoring

import (
	"github.com/adwitiyaio/arka/dependency"
	"github.com/adwitiyaio/arka/secrets"
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
			sm: dm.Get(secrets.DependencySecretsManager).(secrets.Manager),
		}
		nrm.(*newRelicManager).initialize()
	}
	d.Register(DependencyMonitoringManager, nrm)
}
