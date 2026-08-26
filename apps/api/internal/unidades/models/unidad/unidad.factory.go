package unidad

import (
	"github.com/lucsky/cuid"

	"github.com/Sanaruca/condominio/internal/core"
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

func (f *UnidadFactory) Nueva(
	codigo UnidadCodigo,
	estado estadounidad.EstadoDeUnidad,
	titular_primario sujeto.Titular,
	contacto *sujeto.Persona,
	descripcion string,
) (*Unidad, core.Error) {

	if codigo == "" {
		return nil, core.NewValidationError("el codigo es requerido")
	}

	return &Unidad{
		id:               UnidadID(cuid.New()),
		codigo:           codigo,
		estado:           estado,
		deuda:            f.qf.Assemble(0),
		wallet:           f.qf.Assemble(0),
		titular_primario: titular_primario,
		contacto:         contacto,
		descripcion:      descripcion,
		titulares:        core.NewSetFromSlice([]sujeto.Titular{}, func(it sujeto.Titular) sujeto.Titular { return it }),
	}, nil
}

func (f *UnidadFactory) Assemble(
	id string,
	codigo UnidadCodigo,
	estado estadounidad.EstadoDeUnidad,
	deuda,
	wallet int,
	titular_primario sujeto.Titular,
	contacto *sujeto.Persona,
	descripcion string,
	titulares ...sujeto.Titular,
) Unidad {

	_titulares := core.NewSetFromSlice(
		titulares,
		func(it sujeto.Titular) sujeto.Titular { return it },
	)

	return Unidad{
		id:               UnidadID(id),
		codigo:           codigo,
		estado:           estado,
		deuda:            f.qf.Assemble(int64(deuda)),
		contacto:         contacto,
		titular_primario: titular_primario,
		descripcion:      descripcion,
		titulares:        _titulares,
		wallet:           f.qf.Assemble(int64(wallet)),
	}
}
