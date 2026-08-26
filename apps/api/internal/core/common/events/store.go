package events

type PendingEvents struct {
	events []Event
}

func (p *PendingEvents) AddEvent(event Event) {
	p.events = append(p.events, event)
}

func (p *PendingEvents) Dispatch() []Event {
	result := make([]Event, len(p.events))
	copy(result, p.events)
	return result
}

func (p *PendingEvents) Clear() {
	p.events = nil
}
