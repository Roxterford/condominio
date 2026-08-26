package event

import (
	"github.com/Sanaruca/condominio/internal/core/common/events"
)

type CuotaRegistrada struct {
	events.BaseEvent
	ID string `json:"id"`
}

func NewCuotaRegistrada(id, correlationID, causationID string) CuotaRegistrada {
	return CuotaRegistrada{
		BaseEvent: events.NewBaseEvent(id, correlationID, causationID),
		ID:        id,
	}
}

func (c CuotaRegistrada) EventName() string {
	return "cuota.registrada"
}

func (c CuotaRegistrada) Payload() any {
	return struct {
		ID string `json:"id"`
	}{ID: c.ID}
}
