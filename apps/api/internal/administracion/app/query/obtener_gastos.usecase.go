package query

import (
	"fmt"

	"github.com/Sanaruca/condominio/internal/administracion/models/gasto"
	"github.com/Sanaruca/condominio/internal/core"
	"github.com/Sanaruca/condominio/internal/core/common"
	"github.com/Sanaruca/condominio/internal/core/common/filter"
	"github.com/Sanaruca/condominio/internal/core/context"
	"github.com/Sanaruca/condominio/internal/core/usecase"
)

type ObtenerGastosDTO struct {
	common.Paginator
	Filter filter.Filter[gasto.Gasto]
	filter filter.Clause
}

type ObtenerGastos usecase.Handler[context.BaseContext, ObtenerGastosDTO, *common.Paginated[gasto.Gasto]]

type obtenerGastos struct {
	repo gasto.GastoRepository
}

func NewObtenerGastos(repo gasto.GastoRepository) ObtenerGastos {
	if repo == nil {
		panic("repo is nil")
	}
	return &obtenerGastos{
		repo: repo,
	}
}

// Exec implements [ObtenerGastos].
func (uc *obtenerGastos) Exec(
	ctx context.BaseContext,
	input ObtenerGastosDTO,
) (*common.Paginated[gasto.Gasto], core.Error) {

	if err := input.Validate(); err != nil {
		return nil, err
	}

	return uc.repo.ObtenerTodos(ctx, input.filter, input.Paginator)
}

func (o *ObtenerGastosDTO) Validate() core.Error {
	o.Paginator.Sanitize()
	filter, err := o.Filter.Build()

	if err != nil {
		return core.WrapError(err)
	}

	o.filter = filter

	fmt.Println(o.filter)

	return nil
}
