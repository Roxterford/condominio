package deuda

import (
	"context"

	"github.com/Sanaruca/condominio/internal/core"
	"github.com/Sanaruca/condominio/internal/core/common"
	"github.com/Sanaruca/condominio/internal/core/common/filter"
	"github.com/Sanaruca/condominio/internal/unidades/models/unidad"
)

type DeudaRepository interface {
	GetLastDeudaWhereNotPagada(ctx context.Context, unidad unidad.UnidadCodigo) (*Deuda, core.Error)
	ObtenerDeudasDeUnidadPorID(ctx context.Context, unidadID unidad.UnidadID, paginator common.Paginator) (*common.Paginated[Deuda], core.Error)
	ObtenerDeudasDeUnidadPorCodigo(ctx context.Context, unidadCodigo unidad.UnidadCodigo, paginator common.Paginator) (*common.Paginated[Deuda], core.Error)
	Guardar(ctx context.Context, deuda *Deuda) core.Error
	Count(ctx context.Context, filter filter.Clause) (int, core.Error)
}
