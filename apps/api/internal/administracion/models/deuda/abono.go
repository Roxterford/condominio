package deuda

import (
	"time"

	"github.com/Sanaruca/condominio/internal/core/common/quantity"
)

type Abono struct {
	pagoID string
	monto  quantity.Quantity
	fecha  time.Time
}

func (DeudaFactory) AssembleAbono(pagoID string, monto quantity.Quantity, fecha time.Time) *Abono {
	return &Abono{
		pagoID: pagoID,
		monto:  monto,
		fecha:  fecha,
	}
}

func (a *Abono) PagoID() string           { return a.pagoID }
func (a *Abono) Monto() quantity.Quantity { return a.monto }
func (a *Abono) Fecha() time.Time         { return a.fecha }
