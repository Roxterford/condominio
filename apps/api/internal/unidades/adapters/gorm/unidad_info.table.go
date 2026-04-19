package gorm

import "github.com/Sanaruca/condominio/internal/unidades/models/unidad"

type UnidadInfo struct {
	Unidad
	DeudaTotal       int    `gorm:"column:deuda_total"`
	EstadoCuenta     string `gorm:"column:estado_cuenta"`
	CuotasPendientes int    `gorm:"column:cuotas_pendientes"`
}

func (u *UnidadInfo) TableName() string {
	return "unidades_info"
}

func (u *UnidadInfo) ToDomainUnidad(factory *unidad.UnidadFactory) unidad.Unidad {
	return factory.Assemble(
		u.ID,
		u.Codigo,
		u.Estado,
		u.DeudaTotal,
	)
}
