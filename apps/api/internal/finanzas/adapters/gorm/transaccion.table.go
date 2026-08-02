package gorm

import (
	"time"

	"github.com/Sanaruca/condominio/internal/core/common/moneda"
	"github.com/Sanaruca/condominio/internal/core/common/quantity"
	"github.com/Sanaruca/condominio/internal/finanzas/models/transaccion"
	"github.com/Sanaruca/condominio/internal/finanzas/types/metodotransaccion"
	"github.com/Sanaruca/condominio/internal/finanzas/types/roldelmovimiento"
	"github.com/Sanaruca/condominio/internal/finanzas/types/tipodemovimiento"
)

type ITransaccion struct {
	ID            string `gorm:"primaryKey"`
	Fecha         time.Time
	Concepto      string
	MontoTotal    int
	Moneda        moneda.Moneda
	Metodo        metodotransaccion.MetodoDeTransaccion
	Tasa          int
	RegistradoPor string
	CuotaID       *string `gorm:"column:cuota"`
	Registro      time.Time
}

func (ITransaccion) TableName() string {
	return "internal_transacciones"
}

type IMovimiento struct {
	ID            string `gorm:"primaryKey"`
	TransaccionID string
	Tipo          tipodemovimiento.TipoDeMovimiento
	Monto         int
	Rol           roldelmovimiento.RolDelMovimiento
	UnidadCodigo  *string
	ProveedorID   *string
}

func (IMovimiento) TableName() string {
	return "internal_movimientos"
}

type DestinoDePago struct {
	ID         string `gorm:"primaryKey"`
	Movimiento string
	Deuda      string
	Destinado  int
	Fecha      time.Time
}

func (DestinoDePago) TableName() string {
	return "destino_de_pagos"
}

func mapToITransaccion(t *transaccion.TransaccionFinanciera) *ITransaccion {
	return &ITransaccion{
		ID:            t.ID(),
		Fecha:         t.Fecha(),
		Concepto:      t.Concepto(),
		MontoTotal:    int(t.MontoTotal().Value()),
		Moneda:        t.Moneda(),
		Metodo:        t.Metodo(),
		Tasa:          int(t.Tasa().Value()),
		RegistradoPor: t.RegistradoPor(),
		CuotaID:       t.CuotaID(),
		Registro:      t.Registro(),
	}
}

func mapToIMovimientos(t *transaccion.TransaccionFinanciera) []IMovimiento {
	movs := t.Movimientos()
	result := make([]IMovimiento, len(movs))
	for i, m := range movs {
		result[i] = IMovimiento{
			ID:            m.ID(),
			TransaccionID: t.ID(),
			Tipo:          m.Tipo(),
			Monto:         int(m.Monto().Value()),
			Rol:           m.Rol(),
			UnidadCodigo:  m.UnidadCodigo(),
			ProveedorID:   m.ProveedorID(),
		}
	}
	return result
}

func toDomainMovimiento(m IMovimiento, qf *quantity.QuantityFactory) transaccion.Movimiento {
	return *transaccion.AssembleMovimiento(
		m.ID,
		m.Tipo,
		qf.Assemble(int64(m.Monto)),
		m.Rol,
		m.UnidadCodigo,
		m.ProveedorID,
	)
}
