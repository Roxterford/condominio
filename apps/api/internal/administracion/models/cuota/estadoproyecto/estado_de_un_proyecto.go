package estadoproyecto

import "github.com/Sanaruca/condominio/internal/core"

var ErrEstadoInvalido = core.NewInvalidArgumentError("estado de proyecto invalido")

type EstadoDeProyecto string

const (
	BORRADOR EstadoDeProyecto = "BORRADOR"
	ACTIVO   EstadoDeProyecto = "ACTIVO"
	CERRADO  EstadoDeProyecto = "CERRADO"
)

func (e EstadoDeProyecto) String() string {
	return string(e)
}

func (e EstadoDeProyecto) Validate() core.Error {
	switch e {
	case BORRADOR, ACTIVO, CERRADO:
		return nil
	default:
		return ErrEstadoInvalido
	}
}
