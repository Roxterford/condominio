package events

import (
	"context"
	"time"
)

// OutboxEventStoreInterface define el contrato para almacenar y publicar eventos outbox.
// La implementación concreta (GORM, etc.) vive en la capa de infraestructura.
type OutboxEventStoreInterface interface {
	AddEvent(ctx context.Context, event Event) error
	GetPending(ctx context.Context, limit int) ([]*OutboxEventData, error)
	MarkPublished(ctx context.Context, eventID string) error
	MarkFailed(ctx context.Context, eventID string) error
	MoveToDLQ(ctx context.Context, eventID string, errMsg string) error
	CountPending(ctx context.Context) (int64, error)
}

// OutboxEventPublisherInterface define el contrato para el publisher de eventos outbox.
type OutboxEventPublisherInterface interface {
	Start(ctx context.Context) error
	WithBackoff(minBackoff, maxBackoff time.Duration, factor float64) OutboxEventPublisherInterface
}
