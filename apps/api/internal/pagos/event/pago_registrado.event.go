package event

import (
	"time"

	"github.com/Sanaruca/condominio/internal/core/common/quantity"
)

type PagoRegistrado struct {
	ID    string            `json:"id"`
	Monto quantity.Quantity `json:"monto"`
	Fecha time.Time         `json:"fecha"`
}

func NewPagoRegistrado(id string, monto quantity.Quantity, fecha time.Time) PagoRegistrado {
	return PagoRegistrado{
		ID:    id,
		Monto: monto,
		Fecha: fecha,
	}
}

func (p PagoRegistrado) EventName() string {
	return "pago.registrado"
}
