package gormAdapter

import (
	"context"
	"math/rand"
	"time"

	"gorm.io/gorm"

	"github.com/Sanaruca/condominio/internal/core/common/events"
	"github.com/Sanaruca/condominio/internal/core/eventregistry"
	"github.com/Sanaruca/condominio/internal/core/lib/logger"
	"github.com/Sanaruca/condominio/internal/core/lib/metrics"
)

// GormOutboxEventStore implementa events.OutboxEventStoreInterface usando GORM
type GormOutboxEventStore struct {
	db *gorm.DB
}

func NewGormOutboxEventStore(db *gorm.DB) *GormOutboxEventStore {
	return &GormOutboxEventStore{db: db}
}

// DB returns the underlying gorm.DB for advanced queries
func (s *GormOutboxEventStore) DB() *gorm.DB {
	return s.db
}

func (s *GormOutboxEventStore) AddEvent(ctx context.Context, event events.Event) error {
	outboxEvent, err := NewOutboxEvent(event)
	if err != nil {
		return err
	}
	return s.db.WithContext(ctx).Create(outboxEvent).Error
}

func (s *GormOutboxEventStore) GetPending(
	ctx context.Context,
	limit int,
) ([]*events.OutboxEventData, error) {
	var infraEvents []*OutboxEvent
	err := s.db.WithContext(ctx).
		Where("published_at IS NULL").
		Order("created_at ASC").
		Limit(limit).
		Find(&infraEvents).Error
	if err != nil {
		return nil, err
	}

	result := make([]*events.OutboxEventData, len(infraEvents))
	for i, e := range infraEvents {
		result[i] = e.toDomain()
	}
	return result, nil
}

func (s *GormOutboxEventStore) MarkPublished(ctx context.Context, eventID string) error {
	now := time.Now().UTC()
	return s.db.WithContext(ctx).
		Model(&OutboxEvent{}).
		Where("id = ?", eventID).
		Updates(map[string]any{
			"published_at": now,
			"retries":      gorm.Expr("retries + 1"),
		}).Error
}

func (s *GormOutboxEventStore) MarkFailed(ctx context.Context, eventID string) error {
	return s.db.WithContext(ctx).
		Model(&OutboxEvent{}).
		Where("id = ?", eventID).
		Update("retries", gorm.Expr("retries + 1")).Error
}

func (s *GormOutboxEventStore) CountPending(ctx context.Context) (int64, error) {
	var count int64
	err := s.db.WithContext(ctx).
		Model(&OutboxEvent{}).
		Where("published_at IS NULL").
		Count(&count).Error
	return count, err
}

func (s *GormOutboxEventStore) MoveToDLQ(ctx context.Context, eventID string, errMsg string) error {
	var event OutboxEvent
	if err := s.db.WithContext(ctx).Where("id = ?", eventID).First(&event).Error; err != nil {
		return err
	}

	dlqEntry := NewOutboxDLQ(&event, errMsg)
	if err := s.db.WithContext(ctx).Create(dlqEntry).Error; err != nil {
		return err
	}

	return s.db.WithContext(ctx).Where("id = ?", eventID).Delete(&OutboxEvent{}).Error
}

// GormOutboxEventPublisher implementa el publisher con backoff exponencial
type GormOutboxEventPublisher struct {
	store         *GormOutboxEventStore
	publisher     events.EventPublisher
	batchSize     int
	pollInterval  time.Duration
	maxRetries    int
	minBackoff    time.Duration
	maxBackoff    time.Duration
	backoffFactor float64
}

func NewGormOutboxEventPublisher(
	store *GormOutboxEventStore,
	publisher events.EventPublisher,
	batchSize int,
	pollInterval time.Duration,
	maxRetries int,
) *GormOutboxEventPublisher {
	if batchSize <= 0 {
		batchSize = 100
	}
	if pollInterval <= 0 {
		pollInterval = 5 * time.Second
	}
	if maxRetries <= 0 {
		maxRetries = 10
	}
	return &GormOutboxEventPublisher{
		store:         store,
		publisher:     publisher,
		batchSize:     batchSize,
		pollInterval:  pollInterval,
		maxRetries:    maxRetries,
		minBackoff:    1 * time.Second,
		maxBackoff:    60 * time.Second,
		backoffFactor: 2.0,
	}
}

func (p *GormOutboxEventPublisher) WithBackoff(
	minBackoff, maxBackoff time.Duration,
	factor float64,
) *GormOutboxEventPublisher {
	p.minBackoff = minBackoff
	p.maxBackoff = maxBackoff
	p.backoffFactor = factor
	return p
}

func (p *GormOutboxEventPublisher) Start(ctx context.Context) error {
	currentBackoff := p.minBackoff
	ticker := time.NewTicker(p.pollInterval)
	defer ticker.Stop()

	if err := p.updatePendingGauge(ctx); err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			hasEvents, err := p.processBatchWithResult(ctx)
			if err != nil {
				return err
			}

			if err := p.updatePendingGauge(ctx); err != nil {
				return err
			}

			if !hasEvents {
				currentBackoff = p.nextBackoff(currentBackoff)
				ticker.Reset(currentBackoff)
			} else {
				currentBackoff = p.minBackoff
				ticker.Reset(p.pollInterval)
			}
		}
	}
}

func (p *GormOutboxEventPublisher) nextBackoff(current time.Duration) time.Duration {
	next := time.Duration(float64(current) * p.backoffFactor)
	if next > p.maxBackoff {
		next = p.maxBackoff
	}
	jitter := time.Duration(float64(next) * 0.25 * (2*rand.Float64() - 1))
	return next + jitter
}

func (p *GormOutboxEventPublisher) updatePendingGauge(ctx context.Context) error {
	count, err := p.store.CountPending(ctx)
	if err != nil {
		return err
	}
	metrics.OutboxEventsPending.Set(float64(count))
	return nil
}

func (p *GormOutboxEventPublisher) processBatchWithResult(ctx context.Context) (bool, error) {
	start := time.Now()
	domainEvents, err := p.store.GetPending(ctx, p.batchSize)
	if err != nil {
		return false, err
	}

	if len(domainEvents) == 0 {
		return false, nil
	}

	metrics.OutboxEventsPending.Set(float64(len(domainEvents)))

	successCount := 0
	failCount := 0
	dlqCount := 0

	for _, domainEvent := range domainEvents {
		if domainEvent.Retries >= p.maxRetries {
			logger.WarnCtx(ctx, "evento outbox movido a DLQ por max retries",
				"event_id", domainEvent.ID,
				"event_name", domainEvent.EventName,
				"retries", domainEvent.Retries,
			)
			if err := p.store.MoveToDLQ(ctx, domainEvent.ID, "max retries exceeded"); err != nil {
				return false, err
			}
			dlqCount++
			continue
		}

		event, err := eventregistry.DeserializeEvent(domainEvent.EventName, domainEvent.Payload)
		if err != nil {
			failCount++
			logger.ErrorCtx(ctx, err, "fallo al deserializar payload outbox",
				"event_id", domainEvent.ID,
				"event_name", domainEvent.EventName,
			)
			p.store.MarkFailed(ctx, domainEvent.ID)
			continue
		}

		if err := p.publisher.Publish(ctx, event); err != nil {
			failCount++
			metrics.OutboxRetriesTotal.Inc()
			logger.ErrorCtx(ctx, err, "fallo al publicar evento outbox",
				"event_id", domainEvent.ID,
				"event_name", domainEvent.EventName,
				"retries", domainEvent.Retries,
			)
			p.store.MarkFailed(ctx, domainEvent.ID)
			continue
		}

		if err := p.store.MarkPublished(ctx, domainEvent.ID); err != nil {
			return false, err
		}

		successCount++
	}

	metrics.OutboxPublishDuration.Observe(time.Since(start).Seconds())
	metrics.OutboxEventsPublishedTotal.WithLabelValues("success").Add(float64(successCount))
	metrics.OutboxEventsPublishedTotal.WithLabelValues("failed").Add(float64(failCount))
	metrics.OutboxDLQTotal.Add(float64(dlqCount))

	return true, nil
}
