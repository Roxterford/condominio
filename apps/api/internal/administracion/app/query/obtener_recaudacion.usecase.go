package query

import (
	"github.com/Sanaruca/condominio/internal/administracion/models/cuota"
	"github.com/Sanaruca/condominio/internal/core"
	cc "github.com/Sanaruca/condominio/internal/core/context"
	"github.com/Sanaruca/condominio/internal/core/usecase"
)

type ObtenerRecaudacionDTO struct {
	CuotaID string
}

type ObtenerRecaudacion usecase.Handler[cc.BaseContext, ObtenerRecaudacionDTO, *cuota.Recaudacion]

type obtenerRecaudacion struct {
	finder cuota.RecaudacionFinder
}

func NewObtenerRecaudacion(finder cuota.RecaudacionFinder) ObtenerRecaudacion {
	return &obtenerRecaudacion{
		finder: finder,
	}
}

func (u *obtenerRecaudacion) Exec(
	ctx cc.BaseContext,
	input ObtenerRecaudacionDTO,
) (*cuota.Recaudacion, core.Error) {

	r, err := u.finder.ObtenerRecaudacion(ctx, cuota.CuotaID(input.CuotaID))

	if err != nil {
		return nil, err
	}

	return r, nil
}
