package unidad

import (
	"context"

	"github.com/Sanaruca/condominio/internal/core"
	"github.com/Sanaruca/condominio/internal/core/common"
	"github.com/Sanaruca/condominio/internal/core/common/filter"
	"github.com/Sanaruca/condominio/internal/unidades/models/unidad/estadounidad"
)

type UnidadRepository interface {
	Exists(ctx context.Context, unidad UnidadID) (bool, core.Error)
	Obtener(
		ctx context.Context,
		filter filter.Clause,
		paginator common.Paginator,
	) (*common.Paginated[Unidad], core.Error)
	ObtenerTodas(ctx context.Context) ([]string, core.Error)
	ObtenerEstado(ctx context.Context, unidad string) (estadounidad.EstadoDeUnidad, core.Error)
}
