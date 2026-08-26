package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/Sanaruca/condominio/internal/core/common/events"
)

// Compile-time check: IdempotencyStore implements events.IdempotencyStore
var _ events.IdempotencyStore = (*IdempotencyStore)(nil)

// IdempotencyStore provides idempotency checks using Redis
type IdempotencyStore struct {
	client *redis.Client
	ttl    time.Duration
}

func NewIdempotencyStore(client *redis.Client, ttl time.Duration) *IdempotencyStore {
	if ttl <= 0 {
		ttl = 24 * time.Hour
	}
	return &IdempotencyStore{client: client, ttl: ttl}
}

func (s *IdempotencyStore) CheckAndMark(ctx context.Context, key string) (bool, error) {
	result, err := s.client.SetNX(ctx, "idempotency:"+key, "1", s.ttl).Result()
	if err != nil {
		return false, err
	}
	return result, nil
}

func (s *IdempotencyStore) Delete(ctx context.Context, key string) error {
	return s.client.Del(ctx, "idempotency:"+key).Err()
}

// IdempotentDispatcher wraps a dispatcher with idempotency checks
// Works with raw byte payload (for dispatcher.Run)
type IdempotentDispatcher struct {
	idempotency *IdempotencyStore
	keyPrefix   string
	mu          sync.RWMutex
	handlers    map[string]func([]byte) error
}

func NewIdempotentDispatcher(
	idempotency *IdempotencyStore,
	keyPrefix string,
) *IdempotentDispatcher {
	return &IdempotentDispatcher{
		idempotency: idempotency,
		keyPrefix:   keyPrefix,
		handlers:    make(map[string]func([]byte) error),
	}
}

func (d *IdempotentDispatcher) RegisterHandler(name string, handle func([]byte) error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.handlers[name] = handle
}

func (d *IdempotentDispatcher) Run(ctx context.Context, name string, payload []byte) error {
	d.mu.RLock()
	handler, ok := d.handlers[name]
	d.mu.RUnlock()

	if !ok {
		return nil
	}

	eventID := extractEventID(payload)
	if eventID == "" {
		return handler(payload)
	}

	idempotencyKey := d.keyPrefix + ":" + name + ":" + eventID

	acquired, err := d.idempotency.CheckAndMark(ctx, idempotencyKey)
	if err != nil {
		return fmt.Errorf("idempotency check failed: %w", err)
	}

	if !acquired {
		return nil
	}

	if err := handler(payload); err != nil {
		d.idempotency.Delete(ctx, idempotencyKey)
		return err
	}

	return nil
}

func extractEventID(payload []byte) string {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(payload, &raw); err != nil {
		return ""
	}

	if eventIDBytes, ok := raw["event_id"]; ok {
		var eventID string
		if err := json.Unmarshal(eventIDBytes, &eventID); err == nil {
			return eventID
		}
	}

	for _, field := range []string{"id", "EventID", "eventID"} {
		if val, ok := raw[field]; ok {
			var id string
			if err := json.Unmarshal(val, &id); err == nil && id != "" {
				return id
			}
		}
	}

	return ""
}
