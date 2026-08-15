package query

import (
	"github.com/Sanaruca/condominio/internal/core"
	"github.com/Sanaruca/condominio/internal/core/common"
	"github.com/Sanaruca/condominio/internal/core/common/filter"
	cc "github.com/Sanaruca/condominio/internal/core/context"
	"github.com/Sanaruca/condominio/internal/core/usecase"
	"github.com/Sanaruca/condominio/internal/finanzas/models/operacion"
)

type ObtenerGastosDTO struct {
	common.Paginator
	Filter *filter.Filter[operacion.Gasto]
}

type ObtenerGastos usecase.Handler[cc.BaseContext, ObtenerGastosDTO, *common.Paginated[operacion.Gasto]]

type obtenerGastos struct {
	finder operacion.GastoFinder
}

func NewObtenerGastos(finder operacion.GastoFinder) ObtenerGastos {
	if finder == nil {
		panic("finder is nil")
	}
	return &obtenerGastos{finder}
}

func (o *obtenerGastos) Exec(
	ctx cc.BaseContext,
	input ObtenerGastosDTO,
) (*common.Paginated[operacion.Gasto], core.Error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	ftr, e := input.Filter.Build()
	if e != nil {
		return nil, core.WrapError(e)
	}

	return o.finder.Buscar(ctx, ftr, input.Paginator)
}

func (input *ObtenerGastosDTO) Validate() core.Error {
	input.Paginator.Sanitize()
	return nil
}
