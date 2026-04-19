package query

import (
	"github.com/Sanaruca/condominio/internal/core"
	"github.com/Sanaruca/condominio/internal/core/common"
	"github.com/Sanaruca/condominio/internal/core/common/filter"
	cc "github.com/Sanaruca/condominio/internal/core/context"
	"github.com/Sanaruca/condominio/internal/core/usecase"
	"github.com/Sanaruca/condominio/internal/unidades/models/unidad"
)

type ObtenerUnidadesDTO struct {
	common.Paginator
	Filter filter.Filter[unidad.Unidad]

	filter filter.Clause
}

type ObtenerUnidades usecase.Handler[cc.BaseContext, ObtenerUnidadesDTO, *common.Paginated[unidad.Unidad]]

type obtenerUnidades struct {
	repo unidad.UnidadRepository
}

func NewObtenerUnidades(unidadRepository unidad.UnidadRepository) ObtenerUnidades {
	if unidadRepository == nil {
		panic("unidadRepository is nil")
	}
	return &obtenerUnidades{unidadRepository}
}

// Exec implements [ObtenerUnidades].
func (uc *obtenerUnidades) Exec(
	ctx cc.BaseContext,
	input ObtenerUnidadesDTO,
) (*common.Paginated[unidad.Unidad], core.Error) {
	err := input.Validate()
	if err != nil {
		return nil, err
	}

	return uc.repo.Obtener(ctx, input.filter, input.Paginator)
}

func (dto *ObtenerUnidadesDTO) Validate() core.Error {

	if dto != nil {
		dto.Paginator.Sanitize()
		var err error
		clause, err := dto.Filter.Build()
		if err != nil {
			return core.WrapError(err)
		}
		dto.filter = clause
	}

	return nil
}
