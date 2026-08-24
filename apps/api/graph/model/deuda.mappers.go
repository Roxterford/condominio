package model

import (
	"github.com/Sanaruca/condominio/internal/administracion/models/deuda"
	"github.com/Sanaruca/condominio/internal/core/common/filter"
	"github.com/Sanaruca/condominio/internal/unidades/models/sujeto"
)

func (input *DeudaFilter) ToFilter() filter.Filter[deuda.Deuda] {
	return ApplyFilter[deuda.Deuda](input)
}

func DeudaFromDomain(d deuda.Deuda, titular sujeto.Titular) *Deuda {

	abonos := make([]*Abono, len(d.Abonos()))

	for i, a := range d.Abonos() {
		abonos[i] = &Abono{
			Pago:  a.PagoID(),
			Monto: a.Monto().Float(),
			Fecha: a.Fecha(),
		}
	}

	return &Deuda{
		ID:    d.ID(),
		Cuota: d.CuotaID().String(),
		Unidad: &UnidadIdentifiers{
			ID:     d.Unidad().ID().String(),
			Codigo: d.Unidad().Codigo().String(),
		},
		Monto:    d.Monto().Float(),
		Abonos:   abonos,
		Registro: d.Registro(),
		Estado:   d.Estado(),
		Deuda:    d.Deuda().Float(),
		Titular:  titularToGraphQL(titular),
	}
}

func titularToGraphQL(titular sujeto.Titular) *DeudaTitular {
	if titular == nil {
		return nil
	}

	return &DeudaTitular{
		ID:          titular.ID().String(),
		Email:       titular.Email().String(),
		Telefono:    titular.Telefono().String(),
		Cedula:      titular.Cedula().String(),
		DisplayName: titular.DisplayName(),
	}
}
