package events

import "encoding/json"

type Dispatcher struct {
	handlers map[string]func([]byte) error
}

func NewDispatcher() *Dispatcher {
	return &Dispatcher{handlers: make(map[string]func([]byte) error)}
}

func RegisterHandler[T any](d *Dispatcher, name string, handle func(T) error) {
	d.handlers[name] = func(payload []byte) error {
		var event T
		if err := json.Unmarshal(payload, &event); err != nil {
			return err
		}
		return handle(event)
	}
}

func (d *Dispatcher) Run(name string, payload []byte) error {
	if handler, ok := d.handlers[name]; ok {
		return handler(payload)
	}
	return nil
}
