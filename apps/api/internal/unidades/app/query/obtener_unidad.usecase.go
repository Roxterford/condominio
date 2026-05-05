package query

import (
	"github.com/Sanaruca/condominio/internal/core"
	cc "github.com/Sanaruca/condominio/internal/core/context"
	"github.com/Sanaruca/condominio/internal/core/usecase"
	"github.com/Sanaruca/condominio/internal/unidades/models/unidad"
)

// --- DTO y Validación ---

type ObtenerUnidadDTO struct {
	SearchID string
}

func (input ObtenerUnidadDTO) Validate() core.Error {
	if input.SearchID == "" {
		return core.NewValidationError("el campo de búsqueda es requerido")
	}
	return nil
}

// --- Definición de Handlers ---

type ObtenerUnidad usecase.Handler[cc.BaseContext, ObtenerUnidadDTO, *unidad.Unidad]
type ObtenerUnidadPorCodigo usecase.Handler[cc.BaseContext, ObtenerUnidadDTO, *unidad.Unidad]

// --- Implementación Base (Reutilizable) ---

type baseObtenerUnidad struct {
	repo unidad.UnidadRepository
}

func newBaseObtenerUnidad(repo unidad.UnidadRepository) baseObtenerUnidad {
	if repo == nil {
		panic("unidadRepository es requerido (nil)")
	}
	return baseObtenerUnidad{repo: repo}
}

// --- Casos de Uso Específicos ---

// Búsqueda por ID
type obtenerUnidad struct{ baseObtenerUnidad }

func NewObtenerUnidad(repo unidad.UnidadRepository) ObtenerUnidad {
	return &obtenerUnidad{newBaseObtenerUnidad(repo)}
}

func (uc *obtenerUnidad) Exec(ctx cc.BaseContext, input ObtenerUnidadDTO) (*unidad.Unidad, core.Error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	return uc.repo.ObtenerPorID(ctx, unidad.UnidadID(input.SearchID))
}

// Búsqueda por Código
type obtenerUnidadPorCodigo struct{ baseObtenerUnidad }

func NewObtenerUnidadPorCodigo(repo unidad.UnidadRepository) ObtenerUnidadPorCodigo {
	return &obtenerUnidadPorCodigo{newBaseObtenerUnidad(repo)}
}

func (uc *obtenerUnidadPorCodigo) Exec(ctx cc.BaseContext, input ObtenerUnidadDTO) (*unidad.Unidad, core.Error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	// Asumiendo que existe ObtenerPorCodigo en el repositorio
	return uc.repo.ObtenerPorCodigo(ctx, input.SearchID)
}
