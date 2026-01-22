package query

import (
	"strings"

	"github.com/Sanaruca/condominio/internal/administracion"
	"github.com/Sanaruca/condominio/internal/core"
	gormAdapter "github.com/Sanaruca/condominio/internal/core/adapters/gorm"
	"github.com/Sanaruca/condominio/internal/core/common"
	"github.com/Sanaruca/condominio/internal/core/context"
	"github.com/Sanaruca/condominio/internal/core/usecase"
	"gorm.io/gorm"
)

type ObtenerGastosSegunCuotaDTO struct {
	CuotaID string
	common.Paginator
}

type ObtenerGastosSegunCuota usecase.Handler[context.BaseContext, ObtenerGastosSegunCuotaDTO, *common.Paginated[administracion.Gasto]]

func NewObtenerGastosSegunCuota() ObtenerGastosSegunCuota {
	return &obtenerGastosSegunCuota{}
}

type obtenerGastosSegunCuota struct{}

func (obtenerGastosSegunCuota) Exec(
	ctx context.BaseContext,
	input ObtenerGastosSegunCuotaDTO,
) (*common.Paginated[administracion.Gasto], core.Error) {

	if err := input.Validate(); err != nil {
		return nil, err
	}

	total, err := gorm.G[administracion.Gasto](
		ctx.DB,
	).Where("cuota = ?", input.CuotaID).
		Count(ctx.Context(), "*")
	if err != nil {
		return nil, core.WrapError(err)
	}

	gastos, err := gorm.G[administracion.Gasto](ctx.DB).
		Scopes(gormAdapter.GPaginate(input.Paginator)).
		Where("cuota = ?", input.CuotaID).
		Find(ctx.Context())

	if err != nil {
		return nil, core.WrapError(err)
	}

	return common.NewPaginated(gastos, int(total), input.Paginator), nil
}

func (dto *ObtenerGastosSegunCuotaDTO) Validate() core.Error {

	dto.CuotaID = strings.TrimSpace(dto.CuotaID)

	if dto.CuotaID == "" {
		return core.NewInvalidArgumentError("cuota es requerida")
	}

	return nil
}
