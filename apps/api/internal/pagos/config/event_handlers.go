// Deprecated: Config legacy de pagos. Usar transacciones en su lugar.
package config

import (
	"context"

	"github.com/Sanaruca/condominio/internal/core/common/events"
	coreContext "github.com/Sanaruca/condominio/internal/core/context"
	"github.com/Sanaruca/condominio/internal/pagos/app/command"
	"github.com/Sanaruca/condominio/internal/pagos/event"
)

// EventHandlersConfig contiene la configuración de todos los handlers de eventos
type EventHandlersConfig struct {
	Dispatcher  *events.Dispatcher
	AplicarPago command.AplicarPago
}

// NewEventHandlersConfig crea y configura todos los handlers de eventos para pagos
func NewEventHandlersConfig(
	aplicarPago command.AplicarPago,
) *EventHandlersConfig {

	dispatcher := events.NewDispatcher()

	// Registrar handler para el evento pago.registrado
	events.RegisterHandler(
		dispatcher,
		"pago.registrado",
		func(evento event.PagoRegistrado) error {
			// Para eventos del sistema, usamos el wrapper del contexto
			baseCtx, err := coreContext.Wrap(context.Background()).AsBase()
			if err != nil {
				return err
			}

			// Ejecutar el use case para aplicar el pago a la deuda
			_, execErr := aplicarPago.Exec(baseCtx, command.AplicarPagoDTO{
				PagoID: evento.ID,
			})

			return execErr
		},
	)

	return &EventHandlersConfig{
		Dispatcher:  dispatcher,
		AplicarPago: aplicarPago,
	}
}
