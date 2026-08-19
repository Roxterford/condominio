package events

import (
	"time"

	"github.com/google/uuid"
)

type Event interface {
	EventName() string
	EventID() string
	AggregateID() string
	OccurredAt() time.Time
	CorrelationID() string
	CausationID() string
	Payload() any
}

type BaseEvent struct {
	eventID       string    `json:"event_id"`
	aggregateID   string    `json:"aggregate_id"`
	occurredAt    time.Time `json:"occurred_at"`
	correlationID string    `json:"correlation_id"`
	causationID   string    `json:"causation_id"`
}

func NewBaseEvent(aggregateID, correlationID, causationID string) BaseEvent {
	now := time.Now().UTC()
	return BaseEvent{
		eventID:       uuid.New().String(),
		aggregateID:   aggregateID,
		occurredAt:    now,
		correlationID: correlationID,
		causationID:   causationID,
	}
}

func (e BaseEvent) EventID() string       { return e.eventID }
func (e BaseEvent) AggregateID() string   { return e.aggregateID }
func (e BaseEvent) OccurredAt() time.Time { return e.occurredAt }
func (e BaseEvent) CorrelationID() string { return e.correlationID }
func (e BaseEvent) CausationID() string   { return e.causationID }

func (e BaseEvent) WithCorrelationID(id string) BaseEvent {
	e.correlationID = id
	return e
}

func (e BaseEvent) WithCausationID(id string) BaseEvent {
	e.causationID = id
	return e
}

func (e BaseEvent) WithAggregateID(id string) BaseEvent {
	e.aggregateID = id
	return e
}

func (e BaseEvent) Payload() any { return e }
