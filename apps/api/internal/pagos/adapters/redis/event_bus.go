// Deprecated: Event bus legacy de pagos. Usar transacciones en su lugar.
package redis

import (
	"context"
	"encoding/json"

	"github.com/Sanaruca/condominio/internal/core/common/events"
	"github.com/redis/go-redis/v9"
)

type RedisEventBus struct {
	client *redis.Client
	stream string
}

func NewRedisEventBus(client *redis.Client, streamName string) events.EventBus {
	return &RedisEventBus{
		client: client,
		stream: streamName,
	}
}

func (b *RedisEventBus) Publish(ctx context.Context, evento events.Event) error {
	// 1. Convertimos el evento a JSON
	payload, err := json.Marshal(evento)
	if err != nil {
		return err
	}

	// 2. Lo mandamos a Redis Streams
	// Usamos el nombre del evento como una clave dentro del mensaje
	// para que el Dispatcher sepa qué es.
	err = b.client.XAdd(ctx, &redis.XAddArgs{
		Stream: b.stream,
		Values: map[string]any{
			"nombre":  evento.EventName(),
			"payload": payload,
		},
	}).Err()

	return err
}
