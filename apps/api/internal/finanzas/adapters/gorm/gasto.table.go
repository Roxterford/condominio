package gorm

import (
	"time"

	"github.com/Sanaruca/condominio/internal/administracion/models/proveedor"
	"github.com/Sanaruca/condominio/internal/core/common/moneda"
	"github.com/Sanaruca/condominio/internal/core/common/quantity"
	"github.com/Sanaruca/condominio/internal/finanzas/models/operacion"
	"github.com/Sanaruca/condominio/internal/finanzas/types/metodoperacion"
)

type Proveedor struct {
	ID            string
	RIF           string
	Nombre        string
	Email         string
	Telefono      string
	Direccion     string
	Registro      time.Time
	Actualizacion time.Time
}

func (Proveedor) TableName() string {
	return "proveedores"
}

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

	Proveedor Proveedor `gorm:"foreignKey:ProveedorID"`
}

func (Gasto) TableName() string {
	return "gastos"
}

func mapToGasto(
	g Gasto,
	qf *quantity.QuantityFactory,
	pf *proveedor.ProveedorFactory,
) operacion.Gasto {
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

	proveedor := pf.Assemble(
		g.Proveedor.ID,
		g.Proveedor.RIF,
		g.Proveedor.Nombre,
		g.Proveedor.Email,
		g.Proveedor.Telefono,
		&g.Proveedor.Direccion,
		g.Proveedor.Registro,
		g.Proveedor.Actualizacion,
	)

	if g.ProveedorID != nil && *g.ProveedorID != "" {
		return operacion.NewGastoAProveedor(base, *proveedor)
	}

	return base
}
