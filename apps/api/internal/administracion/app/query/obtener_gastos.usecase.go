package query

import (
	"github.com/Sanaruca/condominio/internal/administracion/models/gasto"
	"github.com/Sanaruca/condominio/internal/core"
	"github.com/Sanaruca/condominio/internal/core/common"
	"github.com/Sanaruca/condominio/internal/core/context"

	"github.com/Sanaruca/condominio/internal/core/usecase"
)

type ObtenerGastosDTO struct {
	common.Paginator
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
	input.Paginator.Sanitize()
	return uc.repo.ObtenerTodos(ctx, input.Paginator)
}
