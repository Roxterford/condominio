package redis

import (
	"context"
	"log"
	"time"

	"github.com/Sanaruca/condominio/internal/core/common/events" // Ajusta la ruta a tu core
	"github.com/redis/go-redis/v9"
)

type RedisEventWorker struct {
	client     *redis.Client
	stream     string
	group      string
	consumer   string
	dispatcher *events.Dispatcher
}

func NewRedisEventWorker(
	client *redis.Client,
	stream, group, consumer string,
	dispatcher *events.Dispatcher,
) *RedisEventWorker {
	return &RedisEventWorker{
		client:     client,
		stream:     stream,
		group:      group,
		consumer:   consumer,
		dispatcher: dispatcher,
	}
}

func (w *RedisEventWorker) Start(ctx context.Context) error {
	// 1. Asegurar que el grupo de consumidores existe
	// El ID "0" indica que empezamos a leer desde el inicio del stream si el grupo es nuevo
	err := w.client.XGroupCreateMkStream(ctx, w.stream, w.group, "0").Err()
	if err != nil && err.Error() != "BUSYGROUP Consumer Group name already exists" {
		return err
	}

	log.Printf("Worker iniciado: escuchando stream '%s' como consumidor '%s'", w.stream, w.consumer)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			// 2. Leer mensajes nuevos (ID ">")
			// Bloqueamos por 5 segundos para no saturar la CPU si no hay eventos
			entries, err := w.client.XReadGroup(ctx, &redis.XReadGroupArgs{
				Group:    w.group,
				Consumer: w.consumer,
				Streams:  []string{w.stream, ">"},
				Count:    1,
				Block:    5 * time.Second,
			}).Result()

			if err != nil {
				if err != redis.Nil { // redis.Nil significa que el timeout de Block expiró sin mensajes
					log.Printf("Error leyendo de Redis: %v", err)
				}
				continue
			}

			for _, stream := range entries {
				for _, message := range stream.Messages {
					w.processMessage(ctx, message)
				}
			}
		}
	}
}

func (w *RedisEventWorker) processMessage(ctx context.Context, msg redis.XMessage) {
	// Extraemos los datos tal como los guardó tu RedisEventBus
	eventName, _ := msg.Values["nombre"].(string)
	payloadStr, _ := msg.Values["payload"].(string)
	payload := []byte(payloadStr)

	// 3. Ejecutar el Dispatcher
	err := w.dispatcher.Run(eventName, payload)

	if err != nil {
		log.Printf("Error procesando evento %s (ID: %s): %v", eventName, msg.ID, err)
		// Aquí decides: si no haces XAck, el mensaje queda en el PEL (Pending Entires List)
		// para ser reintentado después por otro worker.
		return
	}

	// 4. Confirmar éxito (ACK)
	w.client.XAck(ctx, w.stream, w.group, msg.ID)
}
