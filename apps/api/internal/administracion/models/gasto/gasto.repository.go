package gasto

import (
	"context"

	"github.com/Sanaruca/condominio/internal/core"
	"github.com/Sanaruca/condominio/internal/core/common"
)

type GastoRepository interface {
	Guardar(ctx context.Context, gasto Gasto) (*GastoID, core.Error)
	ObtenerPorID(ctx context.Context, id GastoID) (*Gasto, core.Error)
	ObtenerTodos(
		ctx context.Context,
		paginator common.Paginator,
	) (*common.Paginated[Gasto], core.Error)
}
