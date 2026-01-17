package villas

import (
	"time"

	"github.com/Sanaruca/condominio/internal/core/errors"
	"github.com/Sanaruca/condominio/internal/villas/types/estadodeuda"
)

var (
	ErrDeudaNoEncontrada = errors.New(errors.NOT_FOUND, "Deuda no encontrada")
	ErrVillaNoEncontrada = errors.New(errors.NOT_FOUND, "Villa no encontrada")
)

type IDeuda struct {
	ID            string
	Cuota         string
	Villa         int
	Registro      time.Time
	Actualizacion time.Time
}

func (_ IDeuda) TableName() string {
	return "internal_deudas"
}

type Deuda struct {
	ID            string
	Villa         int
	Cuota         string
	Monto         int
	Deuda         int
	Estado        estadodeuda.EstadoDeDeuda
	Registro      time.Time
	Actualizacion time.Time
}

func (_ Deuda) TableName() string {
	return "deudas"
}
