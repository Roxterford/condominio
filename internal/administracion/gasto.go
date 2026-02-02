package administracion

import (
	"time"

	"github.com/Sanaruca/condominio/internal/core/errors"
)

var (
	ErrGastoNotFound = errors.New(errors.NOT_FOUND, "Gasto no encontrado")
)

type IGasto struct {
	ID             string
	Proveedor      string
	Cuota          *string
	Monto          int
	Moneda         string
	Tasa           int
	Fecha          time.Time
	Descripcion    *string
	Registro       time.Time
	Actualizacion  time.Time
	RegistradoPor  string
	ActualizadoPor string
}

func (_ IGasto) TableName() string {
	return "internal_gastos"
}

type Gasto struct {
	IGasto
	Total int
}

func (g Gasto) TableName() string {
	return "gastos"
}
