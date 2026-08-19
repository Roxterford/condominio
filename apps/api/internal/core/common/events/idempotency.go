package events

import "context"

// IdempotencyStore define el contrato para verificar y marcar eventos procesados.
// Permite implementaciones con Redis, base de datos, memoria, etc.
type IdempotencyStore interface {
	// CheckAndMark verifica si una key existe y la marca atómicamente.
	// Retorna true si la key no existía (primera vez), false si ya existía.
	CheckAndMark(ctx context.Context, key string) (bool, error)

	// Delete elimina una key del store (para reintentos tras error).
	Delete(ctx context.Context, key string) error
}

// IdempotentEventHandler envuelve un handler de eventos con comprobación de idempotencia.
// prefix: prefijo para la key (ej. "finanzas", "administracion")
// handler: lógica de negocio real del evento
func IdempotentEventHandler(
	store IdempotencyStore,
	prefix string,
	handler func(ctx context.Context, event Event) error,
) func(ctx context.Context, event Event) error {
	return func(ctx context.Context, event Event) error {
		eventID := event.EventID()
		if eventID == "" {
			return handler(ctx, event)
		}

		idempotencyKey := prefix + ":" + event.EventName() + ":" + eventID

		acquired, err := store.CheckAndMark(ctx, idempotencyKey)
		if err != nil {
			return err
		}

		if !acquired {
			return nil
		}

		if err := handler(ctx, event); err != nil {
			store.Delete(ctx, idempotencyKey)
			return err
		}

		return nil
	}
}
