package operacion

import (
	"context"

	"github.com/Sanaruca/condominio/internal/core"
	"github.com/Sanaruca/condominio/internal/core/common"
	"github.com/Sanaruca/condominio/internal/core/common/filter"
)

type OperacionRepository interface {
	Guardar(ctx context.Context, o *Operacion) core.Error
	ObtenerPorID(ctx context.Context, id string) (*Operacion, core.Error)
	Obtener(
		ctx context.Context,
		filter filter.Clause,
		paginator common.Paginator,
	) (*common.Paginated[Operacion], core.Error)
	Count(ctx context.Context, filter filter.Clause) (int, core.Error)
}
