package event

import (
	"time"

	"github.com/Sanaruca/condominio/internal/core/common/events"
	"github.com/Sanaruca/condominio/internal/core/common/quantity"
)

type TransaccionRegistrada struct {
	events.BaseEvent
	ID         string            `json:"id"`
	MontoTotal quantity.Quantity `json:"monto_total"`
	Fecha      time.Time         `json:"fecha"`
}

func NewTransaccionRegistrada(
	id string,
	montoTotal quantity.Quantity,
	fecha time.Time,
	correlationID, causationID string,
) TransaccionRegistrada {
	return TransaccionRegistrada{
		BaseEvent:  events.NewBaseEvent(id, correlationID, causationID),
		ID:         id,
		MontoTotal: montoTotal,
		Fecha:      fecha,
	}
}

func (e TransaccionRegistrada) EventName() string {
	return "transaccion.registrada"
}

func (e TransaccionRegistrada) Payload() any {
	return struct {
		ID         string            `json:"id"`
		MontoTotal quantity.Quantity `json:"monto_total"`
		Fecha      time.Time         `json:"fecha"`
	}{ID: e.ID, MontoTotal: e.MontoTotal, Fecha: e.Fecha}
}
