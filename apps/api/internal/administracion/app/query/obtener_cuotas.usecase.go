package query

import (
	"github.com/Sanaruca/condominio/internal/administracion/models/cuota"
	"github.com/Sanaruca/condominio/internal/core"
	"github.com/Sanaruca/condominio/internal/core/common"
	"github.com/Sanaruca/condominio/internal/core/common/filter"
	"github.com/Sanaruca/condominio/internal/core/context"
	"github.com/Sanaruca/condominio/internal/core/usecase"
)

type ObtenerCuotasDTO struct {
	common.Paginator
	Filter filter.Filter[cuota.CuotaBase]
	filter filter.Clause
}

type ObtenerCuotas usecase.Handler[context.BaseContext, ObtenerCuotasDTO, *common.Paginated[cuota.Cuota]]

type obtenerCuotas struct {
	repo cuota.CuotaRepository
}

func NewObtenerCuotas(repo cuota.CuotaRepository) ObtenerCuotas {
	if repo == nil {
		panic("repo is nil")
	}
	return &obtenerCuotas{
		repo: repo,
	}
}

func (u *obtenerCuotas) Exec(
	ctx context.BaseContext,
	input ObtenerCuotasDTO,
) (*common.Paginated[cuota.Cuota], core.Error) {

	if err := input.Validate(); err != nil {
		return nil, err
	}

	result, err := u.repo.Obtener(ctx, input.filter, input.Paginator)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (dto *ObtenerCuotasDTO) Validate() core.Error {
	if dto == nil {
		return nil
	}

	dto.Paginator.Sanitize()

	filter, err := dto.Filter.Build()

	if err != nil {
		return core.WrapError(err)
	}

	dto.filter = filter

	return nil
}
