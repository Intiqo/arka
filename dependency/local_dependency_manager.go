package dependency

import "github.com/adwitiyaio/arka/logger"

type localDependencyManager struct {
	registry map[string]interface{}
}

func (dm *localDependencyManager) Register(name string, svc interface{}) {
	dm.registry[name] = svc
}

func (dm localDependencyManager) Get(name string) interface{} {
	svc := dm.registry[name]
	if svc == nil {
		logger.Log.Panic().Str("service", name).Msgf("failed to retrieve dependency, %s", name)
	}
	return svc
}
