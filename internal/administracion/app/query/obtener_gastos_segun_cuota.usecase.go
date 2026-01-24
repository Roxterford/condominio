package query

import (
	"errors"
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

	cuota, err := gorm.G[administracion.Cuota](
		ctx.DB,
	).Select("id").
		Where("id = ?", input.CuotaID).
		Take(ctx.Context())

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, administracion.ErrCuotaNoEncontrada
	}
	if err != nil {
		return nil, core.WrapError(err)
	}

	total, err := gorm.G[administracion.Gasto](
		ctx.DB,
	).Where("cuota = ?", cuota.ID).
		Count(ctx.Context(), "*")
	if err != nil {
		return nil, core.WrapError(err)
	}

	gastos, err := gorm.G[administracion.Gasto](ctx.DB).
		Scopes(gormAdapter.GPaginate(input.Paginator)).
		Where("cuota = ?", cuota.ID).
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
