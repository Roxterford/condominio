package model

import (
	"github.com/Sanaruca/condominio/internal/administracion/models/gasto"
	"github.com/Sanaruca/condominio/internal/administracion/models/proveedor"
)

func GastoYProveedorFromDomain(
	gasto gasto.Gasto,
	proveedor proveedor.Proveedor,
) *GastoWithProveedor {

	return &GastoWithProveedor{
		ID:       string(gasto.ID()),
		Concepto: gasto.Concepto(),
		Proveedor: ProveedorFromDomain(proveedor),
		Cuota:         gasto.Cuota(),
		Monto:         (gasto.Monto()).Float(),
		Moneda:        gasto.Moneda(),
		Tasa:          (gasto.Tasa()).Float(),
		Total:         (gasto.Total()).Float(),
		Fecha:         gasto.Fecha(),
		Descripcion:   gasto.Descripcion(),
		Registro:      gasto.Audit().CreatedAt,
		RegistradoPor: gasto.Audit().CreatedBy,
	}
}

func GastoFromDomain(gasto gasto.Gasto) *Gasto {
	return &Gasto{
		ID:        string(gasto.ID()),
		Concepto:  gasto.Concepto(),
		Proveedor: gasto.Proveedor(),
		Cuota:     gasto.Cuota(),
		Moneda: string(gasto.Moneda()),
		Fecha:         gasto.Fecha(),
		Descripcion:   gasto.Descripcion(),
		Registro:      gasto.Audit().CreatedAt,
		RegistradoPor: gasto.Audit().CreatedBy,
	}
}
