package deuda

import "time"

type Abono struct {
	pagoID string
	monto  int
	fecha  time.Time
}

func (DeudaFactory) AssembleAbono(pagoID string, monto int, fecha time.Time) *Abono {
	return &Abono{
		pagoID: pagoID,
		monto:  monto,
		fecha:  fecha,
	}
}

func (a *Abono) PagoID() string   { return a.pagoID }
func (a *Abono) Monto() int       { return a.monto }
func (a *Abono) Fecha() time.Time { return a.fecha }
