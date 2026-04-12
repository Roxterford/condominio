package cuota

import (
	"context"

	"github.com/Sanaruca/condominio/internal/core"
	"github.com/Sanaruca/condominio/internal/core/common"
	"github.com/Sanaruca/condominio/internal/core/common/filter"
)

type CuotaRepository interface {
	Obtener(
		ctx context.Context,
		filter filter.Clause,
		paginator common.Paginator,
	) (*common.Paginated[Cuota], core.Error)
	ObtenerPorID(
		ctx context.Context,
		id CuotaID,
	) (Cuota, core.Error)
	Count(ctx context.Context, filter filter.Clause) (int, core.Error)
}
