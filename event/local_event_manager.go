package event

type localEventManager struct {
	handlers map[string][]Handler
}

func (s *localEventManager) initialize() {
	s.handlers = make(map[string][]Handler)
}

func (m *localEventManager) Register(event string, handler Handler) error {
	_, ok := m.handlers[event]
	if !ok {
		m.handlers[event] = make([]Handler, 0)
	}
	m.handlers[event] = append(m.handlers[event], handler)
	return nil
}

func (m localEventManager) Publish(event string, data interface{}) error {
	for _, handler := range m.handlers[event] {
		go handler(event, data)
	}
	return nil
}
