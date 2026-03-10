package events

type Event interface {
	EventName() string
}

type EventNotifier struct {
	pending_events []Event
}

func (n *EventNotifier) AddEvent(event Event) {
	n.pending_events = append(n.pending_events, event)
}

func (n *EventNotifier) Dispatch() []Event {
	events := n.pending_events
	n.pending_events = nil
	return events
}
