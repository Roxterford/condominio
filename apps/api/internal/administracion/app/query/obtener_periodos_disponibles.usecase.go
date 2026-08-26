package query

import (
	"time"

	"github.com/Sanaruca/condominio/internal/administracion/models/cuota"
	"github.com/Sanaruca/condominio/internal/administracion/models/periododisponible"
	"github.com/Sanaruca/condominio/internal/core"
	"github.com/Sanaruca/condominio/internal/core/common/periodo"
	"github.com/Sanaruca/condominio/internal/core/context"
	"github.com/Sanaruca/condominio/internal/core/envirotment"
	"github.com/Sanaruca/condominio/internal/core/usecase"
)

type ObtenerPeriodosDisponiblesDTO struct{}

type ObtenerPeriodosDisponibles usecase.Handler[
	context.BaseContext,
	ObtenerPeriodosDisponiblesDTO,
	[]periododisponible.PeriodoDisponible,
]

type obtenerPeriodosDisponibles struct {
	cuotas      cuota.CuotaRepository
	calculadora *periododisponible.CalculadoraDePeriodos
}

func NewObtenerPeriodosDisponibles(
	cuotas cuota.CuotaRepository,
	calculadora *periododisponible.CalculadoraDePeriodos,
) ObtenerPeriodosDisponibles {
	if cuotas == nil {
		panic("cuotas is nil")
	}
	if calculadora == nil {
		panic("calculadora is nil")
	}
	return &obtenerPeriodosDisponibles{
		cuotas:      cuotas,
		calculadora: calculadora,
	}
}

func (u *obtenerPeriodosDisponibles) Exec(
	ctx context.BaseContext,
	input ObtenerPeriodosDisponiblesDTO,
) ([]periododisponible.PeriodoDisponible, core.Error) {

	emitidos, err := u.cuotas.ObtenerPeriodosEmitidos(ctx)
	if err != nil {
		return nil, err
	}

	regla, err := periododisponible.Nuevo(envirotment.GetMaxMesesFuturoCuota())
	if err != nil {
		return nil, err
	}

	actual := periodo.DesdeTime(time.Now())

	return u.calculadora.PeriodosDisponibles(actual, regla, emitidos), nil
}
