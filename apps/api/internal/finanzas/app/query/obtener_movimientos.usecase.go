package query

import (
	"github.com/Sanaruca/condominio/internal/core"
	"github.com/Sanaruca/condominio/internal/core/common"
	"github.com/Sanaruca/condominio/internal/core/common/filter"
	cc "github.com/Sanaruca/condominio/internal/core/context"
	"github.com/Sanaruca/condominio/internal/core/usecase"
	"github.com/Sanaruca/condominio/internal/finanzas/models/transaccion"
	"github.com/Sanaruca/condominio/internal/finanzas/types/roldelmovimiento"
	"github.com/Sanaruca/condominio/internal/finanzas/types/tipodemovimiento"
)

type ObtenerMovimientosDTO struct {
	common.Paginator
	Filter *filter.Filter[transaccion.TransaccionFinanciera]
	Tipo   *tipodemovimiento.TipoDeMovimiento
	Rol    *roldelmovimiento.RolDelMovimiento
}

type ObtenerMovimientos usecase.Handler[cc.BaseContext, ObtenerMovimientosDTO, *common.Paginated[transaccion.TransaccionFinanciera]]

type obtenerMovimientos struct {
	repo transaccion.TransaccionRepository
}

func NewObtenerMovimientos(repo transaccion.TransaccionRepository) ObtenerMovimientos {
	if repo == nil {
		panic("repo is nil")
	}
	return &obtenerMovimientos{repo}
}

func (o *obtenerMovimientos) Exec(
	ctx cc.BaseContext,
	input ObtenerMovimientosDTO,
) (*common.Paginated[transaccion.TransaccionFinanciera], core.Error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	ftr, e := input.Filter.Build()
	if e != nil {
		return nil, core.WrapError(e)
	}

	return o.repo.Obtener(ctx, ftr, input.Paginator)
}

func (input ObtenerMovimientosDTO) Validate() core.Error {
	input.Paginator.Sanitize()
	return nil
}
