package model

import (
	"github.com/Sanaruca/condominio/internal/core/common/filter"
	"github.com/Sanaruca/condominio/internal/finanzas/models/operacion"
)

func (f *GastoFilter) ToFilter() filter.Filter[operacion.GastoBase] {

	return applyFilter[operacion.GastoBase](f)
}

func GastoTypeFromDomain(gasto operacion.Gasto) GastoType {

	gasto_proveedor := gasto.AsAProveedor()

	if gasto_proveedor != nil {
		return GastoAProveedor{
			Operacion:     gasto.Operacion(),
			Cuota:         gasto.Cuota(),
			Fecha:         gasto.Fecha(),
			Concepto:      gasto.Concepto(),
			Monto:         gasto.Monto().Float(),
			Moneda:        gasto.Moneda(),
			Total:         gasto.Total().Float(),
			Metodo:        gasto.Metodo(),
			Tasa:          gasto.Tasa().Float(),
			RegistradoPor: gasto.RegistradoPor(),
			Registro:      gasto.Registro(),
			Proveedor:     ProveedorFromDomain(gasto_proveedor.Proveedor()),
			Transaccion:   gasto.Transaccion(),
		}
	}

	return GastoACondominio{
		Operacion:     gasto.Operacion(),
		Cuota:         gasto.Cuota(),
		Fecha:         gasto.Fecha(),
		Concepto:      gasto.Concepto(),
		Monto:         gasto.Monto().Float(),
		Moneda:        gasto.Moneda(),
		Total:         gasto.Total().Float(),
		Metodo:        gasto.Metodo(),
		Tasa:          gasto.Tasa().Float(),
		RegistradoPor: gasto.RegistradoPor(),
		Registro:      gasto.Registro(),
		Transaccion:   gasto.Transaccion(),
	}
}
