// Deprecated: Repository legacy de gasto. Usar internal/transacciones/TransaccionRepository en su lugar.
package gasto

import (
	"context"

	"github.com/Sanaruca/condominio/internal/core"
	"github.com/Sanaruca/condominio/internal/core/common"
	"github.com/Sanaruca/condominio/internal/core/common/filter"
)

type GastoRepository interface {
	Guardar(ctx context.Context, gasto Gasto) (*GastoID, core.Error)
	ObtenerPorID(ctx context.Context, id GastoID) (*Gasto, core.Error)
	ObtenerPorIDs(ctx context.Context, ids []GastoID) ([]Gasto, core.Error)
	Actualizar(ctx context.Context, gasto Gasto) core.Error
	ObtenerTodos(
		ctx context.Context,
		filter filter.Clause,
		paginator common.Paginator,
	) (*common.Paginated[Gasto], core.Error)
}
