package gorm

import (
	"time"

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

func (_ Pago) TableName() string {
	return "pagos"
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
			Destinado: d.Destinado(),
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
		Monto:          p.Monto(),
		Referencia:     p.Referencia(),
		Moneda:         p.Moneda(),
		Tasa:           p.Tasa(),
		Registro:       p.Firma().RegistradoEn(),
		RegistradoPor:  p.Firma().RegistradoPor(),
		Actualizacion:  p.Firma().ActualizadoEn(),
		ActualizadoPor: p.Firma().ActualizadoPor(),
	}
}
