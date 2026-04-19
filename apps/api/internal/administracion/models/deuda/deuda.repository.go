package deuda

import (
	"context"

	"github.com/Sanaruca/condominio/internal/core"
	"github.com/Sanaruca/condominio/internal/core/common/filter"
)

type DeudaRepository interface {
	GetLastDeudaWhereNotPagada(ctx context.Context, unidad string) (*Deuda, core.Error)
	Guardar(ctx context.Context, deuda *Deuda) core.Error
	Count(ctx context.Context, filter filter.Clause) (int, core.Error)
}
