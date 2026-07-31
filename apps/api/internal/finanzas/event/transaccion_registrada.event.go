package event

import (
	"time"

	"github.com/Sanaruca/condominio/internal/core/common/quantity"
)

type TransaccionRegistrada struct {
	ID         string            `json:"id"`
	MontoTotal quantity.Quantity `json:"monto_total"`
	Fecha      time.Time         `json:"fecha"`
}

func NewTransaccionRegistrada(id string, montoTotal quantity.Quantity, fecha time.Time) TransaccionRegistrada {
	return TransaccionRegistrada{
		ID:         id,
		MontoTotal: montoTotal,
		Fecha:      fecha,
	}
}

func (e TransaccionRegistrada) EventName() string {
	return "transaccion.registrada"
}
