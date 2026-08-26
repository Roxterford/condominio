package redis_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"

	redisAdapter "github.com/Sanaruca/condominio/internal/core/adapters/redis"
	"github.com/Sanaruca/condominio/internal/core/common/events"
)

type testEvent struct {
	events.BaseEvent
	ID    string `json:"id"`
	Value int    `json:"value"`
}

func NewTestEvent(id string, value int) testEvent {
	return testEvent{
		BaseEvent: events.NewBaseEvent(id, "test-correlation-id", "test-causation-id"),
		ID:        id,
		Value:     value,
	}
}

func (e testEvent) EventName() string { return "test.event" }

func (e testEvent) Payload() any {
	return struct {
		ID    string `json:"id"`
		Value int    `json:"value"`
	}{ID: e.ID, Value: e.Value}
}

func newTestBus(t *testing.T) (*redisAdapter.RedisEventBus, *redis.Client, *miniredis.Miniredis) {
	t.Helper()

	s, err := miniredis.Run()
	if err != nil {
		t.Fatalf("iniciar miniredis: %v", err)
	}

	client := redis.NewClient(&redis.Options{Addr: s.Addr()})
	bus := redisAdapter.NewRedisEventBus(client)

	t.Cleanup(func() {
		client.Close()
		s.Close()
	})

	return bus, client, s
}

func TestPublishAndConsume(t *testing.T) {
	bus, _, _ := newTestBus(t)
	ctx := context.Background()

	dispatcher := events.NewDispatcher()
	var received testEvent
	done := make(chan struct{})
	events.RegisterHandler(
		dispatcher,
		"test.event",
		func(e testEvent) error {
			received = e
			close(done)
			return nil
		},
	)

	go bus.Consume(ctx, "test-consumer", "test.event", dispatcher)

	event := NewTestEvent("abc-123", 42)
	if err := bus.Publish(ctx, event); err != nil {
		t.Fatalf("publicar evento: %v", err)
	}

	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("el consumidor no procesó el evento a tiempo")
	}

	if received.ID != event.ID || received.Value != event.Value {
		t.Errorf("evento recibido = %+v, se esperaba %+v", received, event)
	}
}

func TestPublishDoesNotLoseEventBeforeConsume(t *testing.T) {
	bus, client, _ := newTestBus(t)
	ctx := context.Background()

	event := NewTestEvent("persisted", 7)
	if err := bus.Publish(ctx, event); err != nil {
		t.Fatalf("publicar evento: %v", err)
	}

	streams, err := client.XRange(ctx, "events:test.event", "-", "+").Result()
	if err != nil {
		t.Fatalf("leer stream: %v", err)
	}
	if len(streams) != 1 {
		t.Fatalf("se esperaba 1 mensaje en el stream, hay %d", len(streams))
	}
	if name, _ := streams[0].Values["name"].(string); name != "test.event" {
		t.Errorf("name del mensaje = %q, se esperaba %q", name, "test.event")
	}
}

func TestFailedHandlerIsNotAcknowledged(t *testing.T) {
	bus, client, _ := newTestBus(t)
	ctx := context.Background()

	dispatcher := events.NewDispatcher()
	var attempts int
	var mu sync.Mutex
	ok := make(chan struct{})
	events.RegisterHandler(
		dispatcher,
		"test.event",
		func(e testEvent) error {
			mu.Lock()
			attempts++
			first := attempts == 1
			mu.Unlock()
			if first {
				return context.DeadlineExceeded
			}
			close(ok)
			return nil
		},
	)

	go bus.Consume(ctx, "test-consumer", "test.event", dispatcher)

	if err := bus.Publish(ctx, NewTestEvent("retry", 1)); err != nil {
		t.Fatalf("publicar evento: %v", err)
	}

	select {
	case <-ok:
	case <-time.After(5 * time.Second):
		t.Fatal("el evento fallido no fue reentregado")
	}

	mu.Lock()
	defer mu.Unlock()
	if attempts < 2 {
		t.Errorf("se esperaban al menos 2 intentos, hubo %d", attempts)
	}

	deadline := time.Now().Add(5 * time.Second)
	for {
		pending, err := client.XPending(ctx, "events:test.event", "events-consumers:test.event").
			Result()
		if err != nil {
			t.Fatalf("leer pendientes: %v", err)
		}
		if pending.Count == 0 {
			break
		}
		if time.Now().After(deadline) {
			t.Errorf("tras éxito no deberían quedar mensajes pendientes, hay %d", pending.Count)
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
}
