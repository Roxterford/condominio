package query

import (
	"github.com/Sanaruca/condominio/internal/administracion/models/cuota"
	"github.com/Sanaruca/condominio/internal/core"
	"github.com/Sanaruca/condominio/internal/core/context"
	"github.com/Sanaruca/condominio/internal/core/usecase"
)

type ObtenerCuotaDTO struct {
	ID string
}

type ObtenerCuota usecase.Handler[context.BaseContext, ObtenerCuotaDTO, cuota.Cuota]

type obtenerCuota struct {
	repo cuota.CuotaRepository
}

func NewObtenerCuota(repo cuota.CuotaRepository) ObtenerCuota {
	if repo == nil {
		panic("repo is nil")
	}
	return &obtenerCuota{repo}
}

// Exec implements [ObtenerCuota].
func (uc *obtenerCuota) Exec(
	ctx context.BaseContext,
	input ObtenerCuotaDTO,
) (cuota.Cuota, core.Error) {

	if err := input.Validate(); err != nil {
		return nil, err
	}

	return uc.repo.ObtenerPorID(ctx, cuota.CuotaID(input.ID))

}

func (dto *ObtenerCuotaDTO) Validate() core.Error {

	if dto.ID == "" {
		return core.NewInvalidArgumentError("id es requerido")
	}

	return nil
}
