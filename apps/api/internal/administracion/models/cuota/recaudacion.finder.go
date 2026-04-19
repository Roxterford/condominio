package cuota

import (
	"context"

	"github.com/Sanaruca/condominio/internal/administracion/types/tipodecuota"
	"github.com/Sanaruca/condominio/internal/core"
	"github.com/Sanaruca/condominio/internal/core/common"
	"github.com/Sanaruca/condominio/internal/core/common/filter"
	"github.com/Sanaruca/condominio/internal/core/common/mes"
	"github.com/Sanaruca/condominio/internal/core/common/quantity"
)

type Recaudacion struct {
	Cuota              CuotaID
	Tipo               tipodecuota.TipoDeCuota
	Mes                mes.Mes
	Anio               int
	MontoCuota         quantity.Quantity
	MontoRecaudado     quantity.Quantity
	MontoPendiente     quantity.Quantity
	MontoEstimado      quantity.Quantity
	Unidades           int
	UnidadesAplicadas  int
	UnidadesSolventes  int
	UnidadesPendientes int
	PagosAsociados     int
}

type RecaudacionFinder interface {
	ObtenerRecaudacion(ctx context.Context, id CuotaID) (*Recaudacion, core.Error)
	Buscar(
		ctx context.Context,
		filter filter.Clause,
	) (*common.Paginated[Recaudacion], core.Error)
}
