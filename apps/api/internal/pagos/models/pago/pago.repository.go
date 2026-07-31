// Deprecated: Este repositorio pertenece al sistema legacy de pagos.
// Usar internal/transacciones/models/transaccion/TransaccionRepository en su lugar.
package pago

import (
	"context"

	"github.com/Sanaruca/condominio/internal/core"
	"github.com/Sanaruca/condominio/internal/core/common"
	"github.com/Sanaruca/condominio/internal/core/common/filter"
)

type PagoRepository interface {
	Obtener(
		ctx context.Context,
		filter filter.Clause,
		paginator common.Paginator,
	) (*common.Paginated[Pago], core.Error)

	// Guardar guarda o actualiza un pago
	Guardar(ctx context.Context, pago *Pago) core.Error
	// GetByID encuentra un pago por su ID
	GetByID(ctx context.Context, id string) (*Pago, core.Error)
	Count(ctx context.Context, filter filter.Clause) (int, core.Error)

	// ObtenerPagosSinDestinos retorna todos los pagos que no tienen destinos asignados.
	ObtenerPagosSinDestinos(ctx context.Context) ([]Pago, core.Error)
}
