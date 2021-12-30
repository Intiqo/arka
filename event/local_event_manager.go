package event

type localEventManager struct {
	handlers map[string][]Handler
}

func (m *localEventManager) Register(event string, handler Handler) {
	_, ok := m.handlers[event]
	if !ok {
		m.handlers[event] = make([]Handler, 0)
	}
	m.handlers[event] = append(m.handlers[event], handler)
}

func (m localEventManager) Publish(event string, data interface{}) {
	for _, handler := range m.handlers[event] {
		go handler(event, data)
	}
}
