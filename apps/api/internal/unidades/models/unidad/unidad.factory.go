package unidad

import (
	"github.com/Sanaruca/condominio/internal/core/common/quantity"
	"github.com/Sanaruca/condominio/internal/unidades/models/sujeto"
	"github.com/Sanaruca/condominio/internal/unidades/models/unidad/estadounidad"
)

type UnidadFactory struct {
	qf *quantity.QuantityFactory
}

func NewUnidadFactory(qf *quantity.QuantityFactory) *UnidadFactory {
	if qf == nil {
		panic("qf is nil")
	}
	return &UnidadFactory{qf}
}

func (f *UnidadFactory) Nueva() *Unidad {
	panic("todo")
}

func (f *UnidadFactory) Assemble(
	id string,
	codigo string,
	estado estadounidad.EstadoDeUnidad,
	deuda int,
	contacto *sujeto.Persona,
) Unidad {
	return Unidad{
		id:       UnidadID(id),
		codigo:   codigo,
		estado:   estado,
		deuda:    f.qf.Assemble(int64(deuda)),
		contacto: contacto,
	}
}
