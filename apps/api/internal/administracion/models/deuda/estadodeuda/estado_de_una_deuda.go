package estadodeuda

import "github.com/Sanaruca/condominio/internal/core"

type EstadoDeDeuda string

const (
	Pendiente EstadoDeDeuda = "PENDIENTE"
	Abonada   EstadoDeDeuda = "ABONADA"
	Saldada   EstadoDeDeuda = "SALDADA"
)

func (e EstadoDeDeuda) Validate() core.Error {
	switch e {
	case Pendiente, Abonada, Saldada:
		return nil
	}

	return core.NewInvalidArgumentError("Estado de deuda %s es invalido", e)
}
