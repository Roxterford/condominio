package unidad

import (
	"github.com/Sanaruca/condominio/internal/core/common/filter"
	"github.com/Sanaruca/condominio/internal/core/errors"
	"github.com/Sanaruca/condominio/internal/unidades/models/unidad/estadounidad"
)

var (
	ErrUnidadNoEncontrada = errors.New(errors.NOT_FOUND, "Unidad no encontrada")
)

type UnidadID string

func (u UnidadID) String() string {
	return string(u)
}

type Unidad struct {
	id     UnidadID
	codigo string
	estado estadounidad.EstadoDeUnidad
	deuda  int
}

func (u *Unidad) ID() string                          { return u.id.String() }
func (u *Unidad) Codigo() string                      { return u.codigo }
func (u *Unidad) Estado() estadounidad.EstadoDeUnidad { return u.estado }
func (u *Unidad) Deuda() int                          { return u.deuda }

func (u *Unidad) PoseeDeuda() bool { return u.deuda > 0 }

func (u Unidad) FilterSpec() filter.Spec {
	return filter.Spec{
		"codigo": filter.TypeString,
	}
}
