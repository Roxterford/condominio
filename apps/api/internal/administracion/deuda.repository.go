package administracion

import (
	"context"

	"github.com/Sanaruca/condominio/internal/core"
)

type DeudaRepository interface {
	GetLastDeudaWhereNotPagada(ctx context.Context, villa int) (*Deuda, core.Error)
}
