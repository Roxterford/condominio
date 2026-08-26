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
	Filter *filter.Filter[operacion.GastoBase]

	filter filter.Clause
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

	return o.finder.Buscar(ctx, input.filter, input.Paginator)
}

func (input *ObtenerGastosDTO) Validate() core.Error {
	input.Paginator.Sanitize()

	if input.Filter != nil {

		ftr, err := input.Filter.Build()

		if err != nil {
			return core.WrapError(err)
		}

		input.filter = ftr
	}

	return nil
}
