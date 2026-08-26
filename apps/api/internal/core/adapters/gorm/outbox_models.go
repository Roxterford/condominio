package gormAdapter

import (
	"encoding/json"
	"time"

	"github.com/Sanaruca/condominio/internal/core/common/events"
)

// OutboxEvent es el modelo de persistencia para la tabla outbox_events.
type OutboxEvent struct {
	ID            string     `gorm:"primaryKey;column:id"`
	EventName     string     `gorm:"column:event_name;not null;index"`
	Payload       []byte     `gorm:"column:payload;not null"`
	CorrelationID string     `gorm:"column:correlation_id;index"`
	CausationID   string     `gorm:"column:causation_id;index"`
	OccurredAt    time.Time  `gorm:"column:occurred_at;not null"`
	CreatedAt     time.Time  `gorm:"column:created_at;not null;autoCreateTime"`
	PublishedAt   *time.Time `gorm:"column:published_at"`
	Retries       int        `gorm:"column:retries;default:0"`
}

func (OutboxEvent) TableName() string {
	return "outbox_events"
}

func NewOutboxEvent(event events.Event) (*OutboxEvent, error) {
	payload, err := json.Marshal(event.Payload())
	if err != nil {
		return nil, err
	}

	return &OutboxEvent{
		ID:            event.EventID(),
		EventName:     event.EventName(),
		Payload:       payload,
		CorrelationID: event.CorrelationID(),
		CausationID:   event.CausationID(),
		OccurredAt:    event.OccurredAt(),
		CreatedAt:     time.Now().UTC(),
		Retries:       0,
	}, nil
}

func (e *OutboxEvent) toDomain() *events.OutboxEventData {
	return &events.OutboxEventData{
		ID:            e.ID,
		EventName:     e.EventName,
		Payload:       e.Payload,
		CorrelationID: e.CorrelationID,
		CausationID:   e.CausationID,
		OccurredAt:    e.OccurredAt,
		CreatedAt:     e.CreatedAt,
		PublishedAt:   e.PublishedAt,
		Retries:       e.Retries,
	}
}

// OutboxDLQ representa una entrada en la dead letter queue.
type OutboxDLQ struct {
	ID            string    `gorm:"primaryKey;column:id"`
	EventName     string    `gorm:"column:event_name;not null;index"`
	Payload       []byte    `gorm:"column:payload;not null"`
	CorrelationID string    `gorm:"column:correlation_id;index"`
	CausationID   string    `gorm:"column:causation_id;index"`
	OccurredAt    time.Time `gorm:"column:occurred_at;not null"`
	CreatedAt     time.Time `gorm:"column:created_at;not null;autoCreateTime"`
	FailedAt      time.Time `gorm:"column:failed_at;not null"`
	ErrorMessage  string    `gorm:"column:error_message"`
	Retries       int       `gorm:"column:retries;default:0"`
}

func (OutboxDLQ) TableName() string {
	return "outbox_dlq"
}

func NewOutboxDLQ(event *OutboxEvent, errMsg string) *OutboxDLQ {
	return &OutboxDLQ{
		ID:            event.ID,
		EventName:     event.EventName,
		Payload:       event.Payload,
		CorrelationID: event.CorrelationID,
		CausationID:   event.CausationID,
		OccurredAt:    event.OccurredAt,
		CreatedAt:     event.CreatedAt,
		FailedAt:      time.Now().UTC(),
		ErrorMessage:  errMsg,
		Retries:       event.Retries,
	}
}
