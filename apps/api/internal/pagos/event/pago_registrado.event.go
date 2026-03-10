package event

import "time"

type PagoRegistrado struct {
	ID    string    `json:"id"`
	Monto int       `json:"monto"`
	Fecha time.Time `json:"fecha"`
}

func NewPagoRegistrado(id string, monto int, fecha time.Time) PagoRegistrado {
	return PagoRegistrado{
		ID:    id,
		Monto: monto,
		Fecha: fecha,
	}
}

func (p PagoRegistrado) EventName() string {
	return "pago.registrado"
}
