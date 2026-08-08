package transaccion

import (
	"context"

	"github.com/Sanaruca/condominio/internal/core"
	"github.com/Sanaruca/condominio/internal/core/common"
	"github.com/Sanaruca/condominio/internal/core/common/filter"
	"github.com/Sanaruca/condominio/internal/finanzas/types/tipodemovimiento"
)

type TransaccionRepository interface {
	Guardar(ctx context.Context, tx *TransaccionFinanciera) core.Error
	ObtenerPorID(ctx context.Context, id string) (*TransaccionFinanciera, core.Error)
	Obtener(
		ctx context.Context,
		filter filter.Clause,
		paginator common.Paginator,
		// TODO: cambiar, no debemos retornar solo los movimientos del tipo espeificado
		tipo *tipodemovimiento.TipoDeMovimiento,
	) (*common.Paginated[TransaccionFinanciera], core.Error)
	Count(ctx context.Context, filter filter.Clause) (int, core.Error)
}
