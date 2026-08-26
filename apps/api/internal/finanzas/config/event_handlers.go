package config

import (
	"context"
	"log"

	"github.com/Sanaruca/condominio/internal/core/common/events"
	coreContext "github.com/Sanaruca/condominio/internal/core/context"
	"github.com/Sanaruca/condominio/internal/finanzas/event"
)

// RegisterEventHandlers registra en el dispatcher los handlers de eventos de
// finanzas.
func RegisterEventHandlers(dispatcher *events.Dispatcher, idempotency events.IdempotencyStore) {
	idempotentHandler := events.IdempotentEventHandler(
		idempotency,
		"finanzas",
		func(ctx context.Context, e events.Event) error {
			switch evt := e.(type) {
			case event.OperacionRegistrada:
				return handleOperacionRegistrada(ctx, evt)
			case event.TransaccionRegistrada:
				return handleTransaccionRegistrada(ctx, evt)
			}
			return nil
		},
	)

	events.RegisterHandler(
		dispatcher,
		new(event.OperacionRegistrada).EventName(),
		func(e event.OperacionRegistrada) error {
			ctx := coreContext.InjectCorrelationID(context.Background(), e.CorrelationID())
			ctx = coreContext.InjectAggregateID(ctx, e.AggregateID())
			return idempotentHandler(ctx, e)
		},
	)
	events.RegisterHandler(
		dispatcher,
		new(event.TransaccionRegistrada).EventName(),
		func(e event.TransaccionRegistrada) error {
			ctx := coreContext.InjectCorrelationID(context.Background(), e.CorrelationID())
			ctx = coreContext.InjectAggregateID(ctx, e.AggregateID())
			return idempotentHandler(ctx, e)
		},
	)
}

func handleOperacionRegistrada(ctx context.Context, e event.OperacionRegistrada) error {
	log.Printf("evento recibido: %s (operación %s)", e.EventName(), e.ID)
	return nil
}

func handleTransaccionRegistrada(ctx context.Context, e event.TransaccionRegistrada) error {
	log.Printf("evento recibido: %s (transacción %s)", e.EventName(), e.ID)
	return nil
}
