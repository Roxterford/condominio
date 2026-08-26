package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/Sanaruca/condominio/internal/core/common/events"
	coreContext "github.com/Sanaruca/condominio/internal/core/context"
	"github.com/Sanaruca/condominio/internal/core/lib/logger"
	"github.com/Sanaruca/condominio/internal/core/lib/metrics"
)

const (
	streamPrefix = "events"
	groupPrefix  = "events-consumers"
	blockTimeout = 5 * time.Second
)

// RedisEventBus implementa el EventBus sobre Redis Streams.
// Usa consumer groups para garantizar que los eventos no se pierdan:
// un mensaje solo se elimina (XACK) cuando el handler lo procesa con éxito.
type RedisEventBus struct {
	client *redis.Client
}

func NewRedisEventBus(client *redis.Client) *RedisEventBus {
	return &RedisEventBus{client: client}
}

// Publish serializa el evento y lo agrega al stream correspondiente a su nombre.
// Redis Streams persiste el mensaje, por lo que no se pierde aunque no haya consumidor.
func (b *RedisEventBus) Publish(ctx context.Context, event events.Event) error {
	start := time.Now()
	eventName := event.EventName()

	payload, err := json.Marshal(event)
	if err != nil {
		metrics.EventsPublishedTotal.WithLabelValues(eventName, "error").Inc()
		return fmt.Errorf("serializar evento %s: %w", eventName, err)
	}

	stream := streamPrefix + ":" + eventName

	_, err = b.client.XAdd(ctx, &redis.XAddArgs{
		Stream: stream,
		Values: map[string]interface{}{
			"name":         eventName,
			"aggregate_id": event.AggregateID(),
			"payload":      payload,
		},
	}).Result()

	metrics.EventsPublishedDuration.WithLabelValues(eventName).Observe(time.Since(start).Seconds())

	if err != nil {
		metrics.EventsPublishedTotal.WithLabelValues(eventName, "error").Inc()
		return fmt.Errorf("publicar evento %s: %w", eventName, err)
	}

	metrics.EventsPublishedTotal.WithLabelValues(eventName, "success").Inc()
	logger.DebugCtx(
		ctx,
		"evento publicado",
		"event_name",
		eventName,
		"aggregate_id",
		event.AggregateID(),
	)

	return nil
}

// Consume escucha el stream del evento y despacha los mensajes al dispatcher.
// Solo confirma (XACK) los mensajes procesados sin error; los fallidos quedan
// pendientes y serán reentregados, garantizando entrega al menos una vez.
func (b *RedisEventBus) Consume(
	ctx context.Context,
	consumer string,
	eventName string,
	dispatcher *events.Dispatcher,
) error {
	stream := streamPrefix + ":" + eventName
	group := groupPrefix + ":" + eventName

	if err := b.ensureGroup(ctx, stream, group); err != nil {
		return err
	}

	logger.InfoCtx(ctx, "Escuchando eventos", "event_name", eventName, "consumer", consumer)

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
		}

		streams, err := b.client.XReadGroup(ctx, &redis.XReadGroupArgs{
			Group:    group,
			Consumer: consumer,
			Streams:  []string{stream, ">"},
			Count:    10,
			Block:    blockTimeout,
		}).Result()

		if err != nil && err != redis.Nil {
			logger.ErrorCtx(ctx, err, "leyendo eventos", "event_name", eventName)
			time.Sleep(time.Second)
			continue
		}

		if err == nil {
			for _, s := range streams {
				for _, msg := range s.Messages {
					b.dispatch(ctx, eventName, dispatcher, msg)
				}
			}
		}

		b.retryPending(ctx, stream, group, consumer, eventName, dispatcher)
	}
}

// dispatch procesa un mensaje y lo confirma solo si el handler termina sin error.
// Si falla, el mensaje permanece en la lista pendiente del grupo y será reentregado
// por retryPending, garantizando entrega al menos una vez.
func (b *RedisEventBus) dispatch(
	ctx context.Context,
	eventName string,
	dispatcher *events.Dispatcher,
	msg redis.XMessage,
) {
	start := time.Now()
	payload, _ := msg.Values["payload"].(string)
	aggregateID, _ := msg.Values["aggregate_id"].(string)

	// Inyectar aggregate_id en el context para que handlers puedan usarlo para ordenamiento
	ctx = coreContext.InjectCorrelationID(ctx, aggregateID)

	err := dispatcher.Run(eventName, []byte(payload))
	metrics.HandlerDuration.WithLabelValues(eventName).Observe(time.Since(start).Seconds())

	if err != nil {
		metrics.EventsConsumedTotal.WithLabelValues(eventName, "error").Inc()
		logger.ErrorCtx(ctx, err, "procesando evento", "event_name", eventName, "msg_id", msg.ID)
		return
	}

	metrics.EventsConsumedTotal.WithLabelValues(eventName, "success").Inc()
	b.client.XAck(ctx, streamPrefix+":"+eventName, groupPrefix+":"+eventName, msg.ID)
	logger.DebugCtx(
		ctx,
		"evento procesado",
		"event_name",
		eventName,
		"aggregate_id",
		aggregateID,
		"msg_id",
		msg.ID,
	)
}

// retryPending reclama mensajes que quedaron pendientes (handlers fallidos)
// y reintenta su procesamiento para que el evento no se pierda.
func (b *RedisEventBus) retryPending(
	ctx context.Context,
	stream, group, consumer, eventName string,
	dispatcher *events.Dispatcher,
) {
	claimedMessages, _, err := b.client.XAutoClaim(ctx, &redis.XAutoClaimArgs{
		Stream:   stream,
		Group:    group,
		Consumer: consumer,
		Start:    "0-0",
		Count:    10,
	}).Result()
	if err != nil {
		logger.ErrorCtx(ctx, err, "reclamando eventos pendientes", "event_name", eventName)
		return
	}

	if len(claimedMessages) > 0 {
		metrics.EventsRetriedTotal.WithLabelValues(eventName).Add(float64(len(claimedMessages)))
	}

	for _, msg := range claimedMessages {
		b.dispatch(ctx, eventName, dispatcher, msg)
	}
}

// ensureGroup crea el consumer group desde el inicio del stream ("0") para que
// los eventos publicados antes de que exista el grupo no se pierdan.
func (b *RedisEventBus) ensureGroup(ctx context.Context, stream, group string) error {
	err := b.client.XGroupCreateMkStream(ctx, stream, group, "0").Err()
	if err != nil && !isBusyGroup(err) {
		return fmt.Errorf("crear grupo %s: %w", group, err)
	}
	return nil
}

func isBusyGroup(err error) bool {
	return err != nil && err.Error() == "BUSYGROUP Consumer Group name already exists"
}

func (b *RedisEventBus) Close() error {
	return b.client.Close()
}
