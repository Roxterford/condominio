package gorm

import (
	"time"

	"github.com/Sanaruca/condominio/internal/core/common/quantity"
	"github.com/Sanaruca/condominio/internal/pagos/models/pago"
	"github.com/Sanaruca/condominio/internal/pagos/types/metododepago"
	"github.com/Sanaruca/condominio/internal/pagos/types/moneda"
)

type IPago struct {
	ID             string
	Unidad         string
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

func (Pago) TableName() string {
	return "pagos"
}

func (p Pago) ToDomain(factory *pago.PagoFactory, qf *quantity.QuantityFactory) *pago.Pago {
	return factory.Assemble(
		p.ID,
		p.Unidad,
		p.Fecha,
		p.Metodo,
		qf.Assemble(int64(p.Monto)),
		p.Moneda,
		qf.Assemble(int64(p.Tasa)),
		p.Referencia,
		p.Registro,
		p.RegistradoPor,
		p.Actualizacion,
		p.ActualizadoPor,
	)
}

type DestinoDePago struct {
	ID        string
	Pago      string
	Deuda     string
	Destinado int
	Fecha     time.Time
}

// mapearDestinos convierte una lista de destinos del dominio a modelos de GORM
func (r *GORMPagoRepository) mapearDestinos(
	destinos []pago.Destino,
	pagoID string,
) []DestinoDePago {
	modelos := make([]DestinoDePago, 0, len(destinos))
	for _, d := range destinos {
		modelos = append(modelos, DestinoDePago{
			ID:        d.ID(), // Usamos el ID generado por el dominio
			Pago:      pagoID,
			Deuda:     d.Deuda(),
			Destinado: int(d.Destinado().Value()),
			Fecha:     d.Fecha(),
		})
	}
	return modelos
}

// mapToIPago convierte la cabecera del pago del dominio al modelo de GORM
func mapToIPago(p *pago.Pago) *IPago {
	return &IPago{
		ID:             p.ID(),
		Unidad:         p.Unidad(),
		Fecha:          p.Fecha(),
		Metodo:         p.Metodo(),
		Monto:          int(p.Monto().Value()),
		Referencia:     p.Referencia(),
		Moneda:         p.Moneda(),
		Tasa:           int(p.Tasa().Value()),
		Registro:       p.Firma().RegistradoEn(),
		RegistradoPor:  p.Firma().RegistradoPor(),
		Actualizacion:  p.Firma().ActualizadoEn(),
		ActualizadoPor: p.Firma().ActualizadoPor(),
	}
}
