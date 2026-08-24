package gorm

import (
	"time"

	"github.com/Sanaruca/condominio/internal/core/common/moneda"
	"github.com/Sanaruca/condominio/internal/core/common/quantity"
	"github.com/Sanaruca/condominio/internal/finanzas/models/operacion"
	"github.com/Sanaruca/condominio/internal/finanzas/types/metodoperacion"
	"github.com/Sanaruca/condominio/internal/finanzas/types/roldestinoperacion"
	"github.com/Sanaruca/condominio/internal/finanzas/types/tipoperacion"
)

// Operacion es la proyeccion GORM de la tabla `internal_operaciones` (base de
// escritura). Fusiona las antiguas tablas internal_transacciones +
// internal_movimientos.
type Operacion struct {
	ID            string `gorm:"primaryKey"`
	Fecha         time.Time
	Concepto      string
	Monto         int
	Moneda        moneda.Moneda
	Metodo        metodoperacion.MetodoDeOperacion
	Tasa          int
	Tipo          tipoperacion.TipoDeOperacion
	Rol           roldestinoperacion.RolDestinoDeOperacion
	Cuota         *string
	UnidadCodigo  *string `gorm:"column:unidad_codigo"`
	ProveedorID   *string `gorm:"column:proveedor"`
	RegistradoPor string
	Registro      time.Time
}

func (Operacion) TableName() string {
	return "internal_operaciones"
}

func mapToOperacion(op *operacion.Operacion) *Operacion {
	return &Operacion{
		ID:            op.ID(),
		Fecha:         op.Fecha(),
		Concepto:      op.Concepto(),
		Monto:         int(op.Monto().Value()),
		Moneda:        op.Moneda(),
		Metodo:        op.Metodo(),
		Tasa:          int(op.Tasa().Value()),
		Tipo:          op.Tipo(),
		Rol:           op.Rol(),
		Cuota:         op.CuotaID(),
		UnidadCodigo:  op.UnidadCodigo(),
		ProveedorID:   op.ProveedorID(),
		RegistradoPor: op.RegistradoPor(),
		Registro:      op.Registro(),
	}
}

func toDomainOperacion(
	o Operacion,
	qf *quantity.QuantityFactory,
	f *operacion.OperacionFactory,
) *operacion.Operacion {
	return f.Assemble(
		o.ID,
		o.Fecha,
		o.Concepto,
		qf.Assemble(int64(o.Monto)),
		o.Moneda,
		o.Metodo,
		qf.Assemble(int64(o.Tasa)),
		o.Tipo,
		o.Rol,
		o.Cuota,
		o.UnidadCodigo,
		o.ProveedorID,
		o.RegistradoPor,
		o.Registro,
	)
}
