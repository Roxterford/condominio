package model

import (
	"github.com/Sanaruca/condominio/internal/administracion/models/deuda"
)

func DeudaFromDomain(d deuda.Deuda) *Deuda {

	abonos := make([]*Abono, len(d.Abonos()))

	for i, a := range d.Abonos() {
		abonos[i] = &Abono{
			Pago:  a.PagoID(),
			Monto: a.Monto().Float(),
			Fecha: a.Fecha(),
		}
	}

	return &Deuda{
		ID:       d.ID(),
		Cuota:    d.CuotaID().String(),
		Unidad:   string(d.Unidad()),
		Monto:    d.Monto().Float(),
		Abonos:   abonos,
		Registro: d.Registro(),
		Estado:   d.Estado(),
		Deuda:    d.Deuda().Float(),
	}
}
