package query

import (
	"context"

	"github.com/Sanaruca/condominio/internal/administracion"
	"github.com/Sanaruca/condominio/internal/core"
	"github.com/Sanaruca/condominio/internal/core/common"
	"github.com/Sanaruca/condominio/internal/core/usecase"
)

type ObtenerCuotasDTO struct {
	common.Paginator
	Filtros *common.FilterSpec[administracion.Cuota]
}

type ObtenerCuotas usecase.WithInOut[ObtenerCuotasDTO, *common.Paginated[administracion.Cuota]]

type obtenerCuotas struct{}

func NewObtenerCuotas() ObtenerCuotas {
	return &obtenerCuotas{}
}

func (o *obtenerCuotas) Exec(
	ctx context.Context,
	input ObtenerCuotasDTO,
) (*common.Paginated[administracion.Cuota], core.Error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	panic("unimplemented")
}

func (dto ObtenerCuotasDTO) Validate() core.Error {

	if dto.Filtros != nil {
		if err := dto.Filtros.Validate(); err != nil {
			return err
		}
		dto.Filtros.Sanitize()
	}

	return nil
}
