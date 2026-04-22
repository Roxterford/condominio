package gorm

import (
	"github.com/Sanaruca/condominio/internal/unidades/models/unidad"
	"github.com/Sanaruca/condominio/internal/unidades/models/unidad/estadounidad"
)

// Deprecated: Use UnidadInfo instead
type Unidad struct {
	ID          string
	Codigo      string
	Estado      estadounidad.EstadoDeUnidad
	Contacto    *string
	Descripcion *string
}

func (u *Unidad) TableName() string {
	return "unidades"
}

// Deprecated
func (u *Unidad) ToDomainUnidad(factory *unidad.UnidadFactory, deuda int) unidad.Unidad {
	return factory.Assemble(
		u.ID,
		u.Codigo,
		u.Estado,
		deuda,
		nil,
	)
}
