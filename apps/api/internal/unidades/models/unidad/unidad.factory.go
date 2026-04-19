package unidad

import "github.com/Sanaruca/condominio/internal/unidades/models/unidad/estadounidad"

type UnidadFactory struct {
}

func NewUnidadFactory() *UnidadFactory {
	return &UnidadFactory{}
}

func (f *UnidadFactory) Nueva() *Unidad {

	panic("todo")

}

func (f *UnidadFactory) Assemble(
	id string,
	codigo string,
	estado estadounidad.EstadoDeUnidad,
	deuda int,

) Unidad {
	return Unidad{
		id:     UnidadID(id),
		codigo: codigo,
		estado: estado,
		deuda:  deuda,
	}
}
