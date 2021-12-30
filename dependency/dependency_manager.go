package dependency

// Manager ... Manages dependencies across the application
type Manager interface {
	// Register ... Registers a service with the dependency manager
	Register(name string, svc interface{})

	// Get ... Gets a service registered with the dependency manager
	Get(name string) interface{}
}

var dm Manager

func GetManager() Manager {
	return dm
}

func init() {
	reg := make(map[string]interface{})
	dm = &localDependencyManager{
		registry: reg,
	}
}
