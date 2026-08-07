package model

import (
	"github.com/Sanaruca/condominio/internal/core/common/filter"
	"github.com/Sanaruca/condominio/internal/finanzas/models/transaccion"
	"github.com/Sanaruca/condominio/internal/finanzas/types/roldelmovimiento"
)

func (input *TransaccionFilter) ToFilter() filter.Filter[transaccion.TransaccionFinanciera] {
	return applyFilter[transaccion.TransaccionFinanciera](input)
}

func TransaccionFromDomain(t transaccion.TransaccionFinanciera) *Transaccion {
	movs := make([]MovimientoType, len(t.Movimientos()))
	for i, m := range t.Movimientos() {
		movs[i] = MovimientoTypeFromDomain(m)
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
		Registro:      t.Registro(),
		Movimientos:   movs,
	}
}

func MovimientoTypeFromDomain(m transaccion.Movimiento) MovimientoType {

	a_condominio := m.AsACondominio()
	a_proveedor := m.AsAProveedor()
	a_unidad := m.AsAUnidad()

	if a_condominio != nil {
		return &MovimientoACondominio{
			ID: m.ID(),
			// Tipo:  a,
			Monto: m.Monto().Float(),
			Cuota: m.CuotaID(),
		}
	}

	if a_proveedor != nil {
		proveedor_id := a_proveedor.Proveedor()
		mov := transaccion.AssembleMovimiento(
			m.ID(),
			m.Tipo(),
			m.Monto(),
			roldelmovimiento.Proveedor,
			nil,
			&proveedor_id,
		)
		return mustProveedor(mov)
	}

	if a_unidad != nil {
		return &MovimientoAUnidad{
			ID:     m.ID(),
			Tipo:   m.Tipo(),
			Monto:  m.Monto().Float(),
			Cuota:  m.CuotaID(),
			Unidad: &Unidad{},
		}
	}

	return nil
}

func mustProveedor(m transaccion.Movimiento) *transaccion.MovimientoAProveedor {
	if p, ok := m.(*transaccion.MovimientoAProveedor); ok {
		return p
	}
	return nil
}
