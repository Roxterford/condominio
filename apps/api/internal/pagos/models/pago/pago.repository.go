package pago

import (
	"context"

	"github.com/Sanaruca/condominio/internal/core"
	"github.com/Sanaruca/condominio/internal/core/common/filter"
)

type PagoRepository interface {
	// Guardar guarda o actualiza un pago
	Guardar(ctx context.Context, pago *Pago) core.Error
	// GetByID encuentra un pago por su ID
	GetByID(ctx context.Context, id string) (*Pago, core.Error)
	Count(ctx context.Context, filter filter.Clause) (int, core.Error)
}
