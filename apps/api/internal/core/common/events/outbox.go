package events

import (
	"time"
)

// OutboxEventData representa los datos de un evento en el outbox.
type OutboxEventData struct {
	ID            string
	EventName     string
	Payload       []byte
	CorrelationID string
	CausationID   string
	OccurredAt    time.Time
	CreatedAt     time.Time
	PublishedAt   *time.Time
	Retries       int
}

// OutboxDLQData representa los datos de un evento en la dead letter queue.
type OutboxDLQData struct {
	ID            string
	EventName     string
	Payload       []byte
	CorrelationID string
	CausationID   string
	OccurredAt    time.Time
	CreatedAt     time.Time
	FailedAt      time.Time
	ErrorMessage  string
	Retries       int
}
