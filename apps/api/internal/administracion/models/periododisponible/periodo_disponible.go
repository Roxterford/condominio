package periododisponible

import (
	"github.com/Sanaruca/condominio/internal/core/common/periodo"
)

type PeriodoDisponible struct {
	periodo          periodo.Periodo
	motivo           string
	esExtraordinario bool
}

func NuevoPeriodoDisponible(
	periodo periodo.Periodo,
	motivo string,
	esExtraordinario bool,
) PeriodoDisponible {
	return PeriodoDisponible{
		periodo:          periodo,
		motivo:           motivo,
		esExtraordinario: esExtraordinario,
	}
}

func (p PeriodoDisponible) Periodo() periodo.Periodo { return p.periodo }
func (p PeriodoDisponible) Motivo() string           { return p.motivo }
func (p PeriodoDisponible) EsExtraordinario() bool   { return p.esExtraordinario }
