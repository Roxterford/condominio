package event

import (
	"time"

	"github.com/Sanaruca/condominio/internal/core/common/events"
	"github.com/Sanaruca/condominio/internal/core/common/quantity"
)

type OperacionRegistrada struct {
	events.BaseEvent
	ID    string            `json:"id"`
	Monto quantity.Quantity `json:"monto"`
	Fecha time.Time         `json:"fecha"`
}

func NewOperacionRegistrada(
	id string,
	monto quantity.Quantity,
	fecha time.Time,
	correlationID, causationID string,
) OperacionRegistrada {
	return OperacionRegistrada{
		BaseEvent: events.NewBaseEvent(id, correlationID, causationID),
		ID:        id,
		Monto:     monto,
		Fecha:     fecha,
	}
}

func (e OperacionRegistrada) EventName() string {
	return "operacion.registrada"
}

func (e OperacionRegistrada) Payload() any {
	return struct {
		ID    string            `json:"id"`
		Monto quantity.Quantity `json:"monto"`
		Fecha time.Time         `json:"fecha"`
	}{ID: e.ID, Monto: e.Monto, Fecha: e.Fecha}
}
