package query

import (
	"github.com/Sanaruca/condominio/internal/administracion"
	"github.com/Sanaruca/condominio/internal/core"
	gormAdapter "github.com/Sanaruca/condominio/internal/core/adapters/gorm"
	"github.com/Sanaruca/condominio/internal/core/common"
	"github.com/Sanaruca/condominio/internal/core/common/filter"
	"github.com/Sanaruca/condominio/internal/core/context"
	"github.com/Sanaruca/condominio/internal/core/usecase"
	"gorm.io/gorm"
)

type ObtenerCuotasDTO struct {
	common.Paginator
	Filter *filter.Filter[administracion.Cuota]
}

type ObtenerCuotas usecase.Handler[context.BaseContext, ObtenerCuotasDTO, *common.Paginated[administracion.Cuota]]

type obtenerCuotas struct{}

func NewObtenerCuotas() ObtenerCuotas {
	return &obtenerCuotas{}
}

func (o *obtenerCuotas) Exec(
	ctx context.BaseContext,
	input ObtenerCuotasDTO,
) (*common.Paginated[administracion.Cuota], core.Error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	var scopes []func(*gorm.Statement)

	if input.Filter != nil {
		println("hay filtros")
		scopes = append(scopes, gormAdapter.GFilter(*input.Filter))
	}

	scopes = append(scopes, gormAdapter.GPaginate(input.Paginator))

	cuotas, err := gorm.G[administracion.Cuota](ctx.DB).Scopes(scopes...).Find(ctx.Context())

	if err != nil {
		println("error: ", err.Error())
		return nil, core.WrapError(err)
	}

	return common.NewPaginated(cuotas, len(cuotas), input.Paginator), nil

}

func (dto ObtenerCuotasDTO) Validate() core.Error {
	if dto.Filter != nil {
		if err := dto.Filter.Validate(); err != nil {
			return err
		}
	}

	return nil
}
