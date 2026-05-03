package unidad

import (
	"github.com/Sanaruca/condominio/internal/core"
	"github.com/Sanaruca/condominio/internal/core/common/filter"
	"github.com/Sanaruca/condominio/internal/core/common/quantity"
	"github.com/Sanaruca/condominio/internal/core/errors"
	"github.com/Sanaruca/condominio/internal/unidades/models/sujeto"
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
	deuda  quantity.Quantity

	// Contacto hace referencia a uno de los titulares de la unidad como contacto
	// principal de la misma. Este campo siempre debería hacer referencia a una
	// persona natural.
	contacto *sujeto.Persona

	titulares core.Set[sujeto.Sujeto]
}

func (u *Unidad) ID() string                          { return u.id.String() }
func (u *Unidad) Codigo() string                      { return u.codigo }
func (u *Unidad) Estado() estadounidad.EstadoDeUnidad { return u.estado }
func (u *Unidad) Deuda() quantity.Quantity            { return u.deuda }
func (u *Unidad) Contacto() *sujeto.Persona {

	if u != nil {
		return u.contacto
	}

	return nil
}
func (u *Unidad) PoseeDeuda() bool { return u.deuda.Value() > 0 }

func (u Unidad) FilterSpec() filter.Spec {
	return filter.Spec{
		"id":     filter.TypeString,
		"codigo": filter.TypeString,
	}
}
