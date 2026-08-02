package model

import (
	"github.com/Sanaruca/condominio/internal/core/common/filter"
	"github.com/Sanaruca/condominio/internal/finanzas/models/transaccion"
)

func (input *TransaccionFilter) ToFilter() filter.Filter[transaccion.TransaccionFinanciera] {
	return applyFilter[transaccion.TransaccionFinanciera](input)
}

func TransaccionFromDomain(t transaccion.TransaccionFinanciera) *Transaccion {
	movs := make([]*Movimiento, len(t.Movimientos()))
	for i, m := range t.Movimientos() {
		movs[i] = MovimientoFromDomain(m)
	}
	return &Transaccion{
		ID:            t.ID(),
		Fecha:         t.Fecha(),
		Concepto:      t.Concepto(),
		MontoTotal:    t.MontoTotal().Float(),
		Moneda:        t.Moneda(),
		Metodo:        t.Metodo(),
		Tasa:          t.Tasa().Float(),
		RegistradoPor: t.RegistradoPor(),
		CuotaID:       t.CuotaID(),
		Registro:      t.Registro(),
		Movimientos:   movs,
	}
}

func MovimientoFromDomain(m transaccion.Movimiento) *Movimiento {
	return &Movimiento{
		ID:           m.ID(),
		Tipo:         TipoDeMovimiento(m.Tipo()),
		Monto:        m.Monto().Float(),
		Rol:          RolDelMovimiento(m.Rol()),
		UnidadCodigo: m.UnidadCodigo(),
		ProveedorID:  m.ProveedorID(),
	}
}
