package query

import (
	"github.com/Sanaruca/condominio/internal/administracion/models/deuda"
	"github.com/Sanaruca/condominio/internal/core"
	"github.com/Sanaruca/condominio/internal/core/common"
	"github.com/Sanaruca/condominio/internal/core/common/filter"
	cc "github.com/Sanaruca/condominio/internal/core/context"
	"github.com/Sanaruca/condominio/internal/core/usecase"
)

type ObtenerDeudoresDTO struct {
	common.Paginator
	Filter *filter.Filter[deuda.Deuda]

	filter filter.Clause
}

type ObtenerDeudores usecase.Handler[cc.BaseContext, ObtenerDeudoresDTO, *common.Paginated[deuda.DeudaConTitular]]

type obtenerDeudores struct {
	deudas deuda.DeudaRepository
}

func NewObtenerDeudores(deudaRepository deuda.DeudaRepository) ObtenerDeudores {

	if deudaRepository == nil {
		panic("deudaRepository is nil")
	}

	return &obtenerDeudores{deudaRepository}
}

// Exec implements [ObtenerDeudores].
func (uc *obtenerDeudores) Exec(
	ctx cc.BaseContext,
	input ObtenerDeudoresDTO,
) (*common.Paginated[deuda.DeudaConTitular], core.Error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	return uc.deudas.ObtenerConTitular(ctx, input.filter, input.Paginator)
}

func (input *ObtenerDeudoresDTO) Validate() core.Error {

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
