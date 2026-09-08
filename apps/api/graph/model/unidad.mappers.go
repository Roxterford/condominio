package model

import (
	"github.com/Sanaruca/condominio/internal/core/common/filter"
	"github.com/Sanaruca/condominio/internal/unidades/models/unidad"
)

func (input *UnidadFilter) ToFilter() filter.Filter[unidad.Unidad] {
	return ApplyFilter[unidad.Unidad](input)
}

func UnidadFromDomain(unidad unidad.Unidad) *Unidad {

	contacto := PersonaFromDomain(unidad.Contacto())
	titular := TitularFromDomain(unidad.TitularPrimario())

	return &Unidad{
		ID:              unidad.ID().String(),
		Codigo:          string(unidad.Codigo()),
		Estado:          unidad.Estado(),
		TitularPrimario: titular,
		Contacto:        contacto,
		Deuda:           unidad.Deuda().Float(),
		Wallet:          unidad.Wallet().Float(),
		Titulares:       []Titular{},
	}
}
