package events

import (
	"encoding/json"
	"sync"
)

type Dispatcher struct {
	mu       sync.RWMutex
	handlers map[string]func([]byte) error
}

func NewDispatcher() *Dispatcher {
	return &Dispatcher{handlers: make(map[string]func([]byte) error)}
}

func (d *Dispatcher) RegisterHandler(name string, handle func([]byte) error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.handlers[name] = handle
}

func RegisterHandler[T Event](d *Dispatcher, name string, handle func(T) error) {
	d.RegisterHandler(name, func(payload []byte) error {
		var event T
		if err := json.Unmarshal(payload, &event); err != nil {
			return err
		}
		return handle(event)
	})
}

func (d *Dispatcher) Run(name string, payload []byte) error {
	d.mu.RLock()
	handler, ok := d.handlers[name]
	d.mu.RUnlock()
	if ok {
		return handler(payload)
	}
	return nil
}
