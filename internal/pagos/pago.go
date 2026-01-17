package pagos

import (
	"time"

	"github.com/Sanaruca/condominio/internal/core/errors"
	"github.com/Sanaruca/condominio/internal/pagos/types/metododepago"
	"github.com/Sanaruca/condominio/internal/pagos/types/moneda"
)

var (
	ErrPagoNoEncontrado = errors.New(errors.NOT_FOUND, "Pago no encontrado")
)

type IPago struct {
	ID             string
	Villa          int
	Fecha          time.Time
	Metodo         metododepago.MetodoDePago
	Monto          int
	Referencia     *string
	Moneda         moneda.Moneda
	Tasa           int
	Registro       time.Time
	RegistradoPor  string
	Actualizacion  time.Time
	ActualizadoPor string
}

func (_ IPago) TableName() string {
	return "internal_pagos"
}

type Pago struct {
	IPago
	Total     int
	Destinado int
	Cuenta    int
}

func (_ Pago) TableName() string {
	return "pagos"
}
