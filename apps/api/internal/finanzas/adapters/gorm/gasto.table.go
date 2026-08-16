package gorm

import (
	"time"

	"github.com/Sanaruca/condominio/internal/core/common/moneda"
	"github.com/Sanaruca/condominio/internal/core/common/quantity"
	"github.com/Sanaruca/condominio/internal/finanzas/models/operacion"
	"github.com/Sanaruca/condominio/internal/finanzas/types/metodoperacion"
)

// Readonly / View
type Gasto struct {
	Operacion     string
	Transaccion   *string
	Cuota         *string
	Fecha         time.Time
	Concepto      string
	Monto         int
	Moneda        moneda.Moneda
	Metodo        metodoperacion.MetodoDeOperacion
	Tasa          int
	UnidadCodigo  *string `gorm:"column:unidad_codigo"`
	ProveedorID   *string `gorm:"column:proveedor"`
	RegistradoPor string
	Registro      time.Time
}

func (Gasto) TableName() string {
	return "gastos"
}

func mapToGasto(g Gasto, qf *quantity.QuantityFactory) operacion.Gasto {
	base := operacion.NewGastoBase(
		g.Operacion,
		g.Cuota,
		g.Fecha,
		g.Concepto,
		qf.Assemble(int64(g.Monto)),
		g.Moneda,
		g.Metodo,
		qf.Assemble(int64(g.Tasa)),
		g.RegistradoPor,
		g.Registro,
	)

	if g.ProveedorID != nil && *g.ProveedorID != "" {
		return operacion.NewGastoAProveedor(base, *g.ProveedorID)
	}

	return base
}
