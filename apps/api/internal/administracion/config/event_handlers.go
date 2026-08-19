package config

import (
	"context"
	"log"

	"github.com/Sanaruca/condominio/internal/administracion/app/command"
	"github.com/Sanaruca/condominio/internal/administracion/event"
	"github.com/Sanaruca/condominio/internal/core/common/events"
	coreContext "github.com/Sanaruca/condominio/internal/core/context"
)

// RegisterEventHandlers registra en el dispatcher los handlers de eventos de
// administración.
func RegisterEventHandlers(
	dispatcher *events.Dispatcher,
	aplicarCuota command.AplicarCuota,
	idempotency events.IdempotencyStore,
) {
	idempotentHandler := events.IdempotentEventHandler(
		idempotency,
		"administracion",
		func(ctx context.Context, e events.Event) error {
			if evt, ok := e.(event.CuotaRegistrada); ok {
				return handleCuotaRegistrada(ctx, evt, aplicarCuota)
			}
			return nil
		},
	)

	events.RegisterHandler(
		dispatcher,
		event.CuotaRegistrada{}.EventName(),
		func(e event.CuotaRegistrada) error {
			ctx := coreContext.InjectCorrelationID(context.Background(), e.CorrelationID())
			ctx = coreContext.InjectAggregateID(ctx, e.AggregateID())
			return idempotentHandler(ctx, e)
		},
	)
}

func handleCuotaRegistrada(
	ctx context.Context,
	e event.CuotaRegistrada,
	aplicarCuota command.AplicarCuota,
) error {
	log.Printf("evento recibido: %s (cuota %s)", e.EventName(), e.ID)

	baseCtx, err := coreContext.Wrap(ctx).AsBase()
	if err != nil {
		return err
	}

	_, execErr := aplicarCuota.Exec(
		baseCtx,
		command.AplicarCuotaDTO{CuotaID: e.ID},
	)

	return execErr
}
