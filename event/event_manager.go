package event

type Handler func(event string, data interface{})

var em Manager

type Manager interface {
	// Register ... Register an event handler
	Register(event string, handler Handler)

	// Publish ... Publish an event
	Publish(event string, data interface{})
}

func GetManager() Manager {
	return em
}

func init() {
	h := make(map[string][]Handler)
	em = &localEventManager{
		handlers: h,
	}
}
