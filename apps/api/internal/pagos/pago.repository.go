package pagos

import (
	"context"

	"github.com/Sanaruca/condominio/internal/core"
)

type PagoRepository interface {
	// Guardar guarda o actualiza un pago
	Guardar(ctx context.Context, pago *Pago) core.Error
	// GetByID encuentra un pago por su ID
	GetByID(ctx context.Context, id string) (*Pago, core.Error)
}
