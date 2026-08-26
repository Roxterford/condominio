package query

import (
	"github.com/Sanaruca/condominio/internal/core"
	"github.com/Sanaruca/condominio/internal/core/common"
	"github.com/Sanaruca/condominio/internal/core/common/filter"
	cc "github.com/Sanaruca/condominio/internal/core/context"
	"github.com/Sanaruca/condominio/internal/core/usecase"
	"github.com/Sanaruca/condominio/internal/finanzas/models/operacion"
)

type ObtenerMovimientosDTO struct {
	common.Paginator
	Filter *filter.Filter[operacion.Operacion]
}

type ObtenerMovimientos usecase.Handler[cc.BaseContext, ObtenerMovimientosDTO, *common.Paginated[operacion.Operacion]]

type obtenerMovimientos struct {
	repo operacion.OperacionRepository
}

func NewObtenerMovimientos(repo operacion.OperacionRepository) ObtenerMovimientos {
	if repo == nil {
		panic("repo is nil")
	}
	return &obtenerMovimientos{repo}
}

func (o *obtenerMovimientos) Exec(
	ctx cc.BaseContext,
	input ObtenerMovimientosDTO,
) (*common.Paginated[operacion.Operacion], core.Error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	ftr, e := input.Filter.Build()
	if e != nil {
		return nil, core.WrapError(e)
	}

	return o.repo.Obtener(ctx, ftr, input.Paginator)
}

func (input *ObtenerMovimientosDTO) Validate() core.Error {
	input.Paginator.Sanitize()
	return nil
}
