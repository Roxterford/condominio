package query

import (
	"github.com/Sanaruca/condominio/internal/core"
	cc "github.com/Sanaruca/condominio/internal/core/context"
	"github.com/Sanaruca/condominio/internal/core/usecase"
	"github.com/Sanaruca/condominio/internal/unidades/models/unidad"
)

type ObtenerUnidadDTO struct {
	ID string
}

type ObtenerUnidad usecase.Handler[cc.BaseContext, ObtenerUnidadDTO, *unidad.Unidad]

type obtenerUnidad struct {
	repo unidad.UnidadRepository
}

func NewObtenerUnidad(unidadRepository unidad.UnidadRepository) ObtenerUnidad {
	if unidadRepository == nil {
		panic("unidadRepository is nil")
	}
	return &obtenerUnidad{unidadRepository}
}

func (uc *obtenerUnidad) Exec(
	ctx cc.BaseContext,
	input ObtenerUnidadDTO,
) (*unidad.Unidad, core.Error) {

	err := input.Validate()

	if err != nil {
		return nil, err
	}

	return uc.repo.ObtenerPorID(ctx, unidad.UnidadID(input.ID))
}

func (input ObtenerUnidadDTO) Validate() core.Error {

	if input.ID == "" {
		return core.NewValidationError("id es requerido")
	}

	return nil
}
