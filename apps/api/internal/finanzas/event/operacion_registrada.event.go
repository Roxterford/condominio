package event

import (
	"time"

	"github.com/Sanaruca/condominio/internal/core/common/quantity"
)

type OperacionRegistrada struct {
	ID    string            `json:"id"`
	Monto quantity.Quantity `json:"monto"`
	Fecha time.Time         `json:"fecha"`
}

func NewOperacionRegistrada(
	id string,
	monto quantity.Quantity,
	fecha time.Time,
) OperacionRegistrada {
	return OperacionRegistrada{
		ID:    id,
		Monto: monto,
		Fecha: fecha,
	}
}

func (e OperacionRegistrada) EventName() string {
	return "operacion.registrada"
}
