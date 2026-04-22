package gorm

import (
	"github.com/Sanaruca/condominio/internal/unidades/models/sujeto"
	"github.com/Sanaruca/condominio/internal/unidades/models/unidad"
)

type UnidadInfo struct {
	Unidad
	DeudaTotal       int    `gorm:"column:deuda_total"`
	EstadoCuenta     string `gorm:"column:estado_cuenta"`
	CuotasPendientes int    `gorm:"column:cuotas_pendientes"`

	PersonaContacto *Sujeto
}

func (u *UnidadInfo) TableName() string {
	return "unidades_info"
}

func (u *UnidadInfo) ToDomainUnidad(
	unidadFactory *unidad.UnidadFactory,
	sujetoFactory *sujeto.SujetoFactory,
) unidad.Unidad {

	var persona *sujeto.Persona
	if u.PersonaContacto != nil {

		var nombres, apellidos string
		if u.PersonaContacto.Nombres != nil {
			nombres = *u.PersonaContacto.Nombres
		}
		if u.PersonaContacto.Apellidos != nil {
			apellidos = *u.PersonaContacto.Apellidos
		}

		persona = sujetoFactory.AssemblePersona(
			u.PersonaContacto.ID,
			u.PersonaContacto.DocumentoIdentidad,
			nombres,
			apellidos,
			u.PersonaContacto.Email,
			u.PersonaContacto.Telefono,
		)
	}

	return unidadFactory.Assemble(
		u.ID,
		u.Codigo,
		u.Estado,
		u.DeudaTotal,
		persona,
	)
}
