package query

import (
	"github.com/Sanaruca/condominio/internal/administracion/models/deuda"
	"github.com/Sanaruca/condominio/internal/core"
	"github.com/Sanaruca/condominio/internal/core/common"
	cc "github.com/Sanaruca/condominio/internal/core/context"
	"github.com/Sanaruca/condominio/internal/core/usecase"
	"github.com/Sanaruca/condominio/internal/unidades/models/unidad"
)

type ObtenerDeudasDeUnaUnidadDTO struct {
	SearchID string
	common.Paginator
}

func (input ObtenerDeudasDeUnaUnidadDTO) Validate() core.Error {

	input.Paginator.Sanitize()

	if input.SearchID == "" {
		return core.NewValidationError("el campo de búsqueda es requerido")
	}

	return nil

}

// --- Definición de Handlers ---

type ObtenerDeudasDeUnaUnidadPorID usecase.Handler[cc.BaseContext, ObtenerDeudasDeUnaUnidadDTO, *common.Paginated[deuda.Deuda]]

type ObtenerDeudasDeUnaUnidadPorCodigo usecase.Handler[cc.BaseContext, ObtenerDeudasDeUnaUnidadDTO, *common.Paginated[deuda.Deuda]]

// --- Implementación Base (Reutilizable) ---

type baseObtenerDeudasDeUnaUnidad struct {
	repo deuda.DeudaRepository
}

func newBaseObtenerDeudasDeUnaUnidad(repo deuda.DeudaRepository) baseObtenerDeudasDeUnaUnidad {
	if repo == nil {
		panic("deudaRepository es requerido (nil)")
	}
	return baseObtenerDeudasDeUnaUnidad{repo: repo}
}

// --- Casos de Uso Específicos ---

// Búsqueda por ID
type obtenerDeudasDeUnaUnidadPorID struct{ baseObtenerDeudasDeUnaUnidad }

func NewObtenerDeudasDeUnaUnidadPorID(repo deuda.DeudaRepository) ObtenerDeudasDeUnaUnidadPorID {
	return &obtenerDeudasDeUnaUnidadPorID{newBaseObtenerDeudasDeUnaUnidad(repo)}
}

func (uc *obtenerDeudasDeUnaUnidadPorID) Exec(
	ctx cc.BaseContext,
	input ObtenerDeudasDeUnaUnidadDTO,
) (*common.Paginated[deuda.Deuda], core.Error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	return uc.repo.ObtenerDeudasDeUnidadPorID(ctx, unidad.UnidadID(input.SearchID), input.Paginator)
}

// Búsqueda por Código
type obtenerDeudasDeUnaUnidadPorCodigo struct{ baseObtenerDeudasDeUnaUnidad }

func NewObtenerDeudasDeUnaUnidadPorCodigo(
	repo deuda.DeudaRepository,
) ObtenerDeudasDeUnaUnidadPorCodigo {
	return &obtenerDeudasDeUnaUnidadPorCodigo{newBaseObtenerDeudasDeUnaUnidad(repo)}
}

func (uc *baseObtenerDeudasDeUnaUnidad) Exec(
	ctx cc.BaseContext,
	input ObtenerDeudasDeUnaUnidadDTO,
) (*common.Paginated[deuda.Deuda], core.Error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	return uc.repo.ObtenerDeudasDeUnidadPorCodigo(ctx, input.SearchID, input.Paginator)
}
