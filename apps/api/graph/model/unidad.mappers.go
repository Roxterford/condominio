package model

import (
	"github.com/Sanaruca/condominio/internal/core/common/filter"
	"github.com/Sanaruca/condominio/internal/unidades/models/unidad"
)

func (input *UnidadFilter) ToFilter() filter.Filter[unidad.Unidad] {
	return ApplyFilter[unidad.Unidad](input)
}

func UnidadesTotalesFromDomain(estadisticas unidad.Estadisticas) *UnidadesTotales {

	return &UnidadesTotales{
		TotalUnidades:         int32(estadisticas.TotalUnidades),
		UnidadesActivas:       int32(estadisticas.UnidadesActivas),
		UnidadesInhabitadas:   int32(estadisticas.UnidadesInhabitadas),
		UnidadesExentas:       int32(estadisticas.UnidadesExentas),
		UnidadesEnLitigio:     int32(estadisticas.UnidadesEnLitigio),
		UnidadesSuspendidas:   int32(estadisticas.UnidadesSuspendidas),
		UnidadesPreventa:      int32(estadisticas.UnidadesPreventa),
		UnidadesConPendientes: int32(estadisticas.UnidadesConPendientes),
		UnidadesSolventes:     int32(estadisticas.UnidadesSolventes),
		TotalPendiente:        estadisticas.TotalPendiente.Float(),
		TotalAsignado:         estadisticas.TotalAsignado.Float(),
	}

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
