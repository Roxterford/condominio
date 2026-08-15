package gorm

import (
	"time"

	"github.com/Sanaruca/condominio/internal/core/common/moneda"
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
