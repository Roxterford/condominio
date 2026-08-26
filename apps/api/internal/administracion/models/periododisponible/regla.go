package periododisponible

import "github.com/Sanaruca/condominio/internal/core"

const DefaultMaximoMesesFuturo = 1

type ReglaDeEmision struct {
	maximoMesesFuturo int
}

func Nuevo(maximoMesesFuturo int) (ReglaDeEmision, core.Error) {
	if maximoMesesFuturo < 0 {
		return ReglaDeEmision{}, core.NewValidationError(
			"el máximo de meses futuros no puede ser negativo",
		)
	}
	return ReglaDeEmision{maximoMesesFuturo: maximoMesesFuturo}, nil
}

func (r ReglaDeEmision) MaximoMesesFuturo() int { return r.maximoMesesFuturo }
