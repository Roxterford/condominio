package operacion

import (
	"context"
	"time"

	"github.com/Sanaruca/condominio/internal/core"
	"github.com/Sanaruca/condominio/internal/core/common"
	"github.com/Sanaruca/condominio/internal/core/common/filter"
	"github.com/Sanaruca/condominio/internal/core/common/moneda"
	"github.com/Sanaruca/condominio/internal/core/common/quantity"
	"github.com/Sanaruca/condominio/internal/finanzas/types/metodoperacion"
)

// GastoBase es un read model (proyección) de una operación DEBITO junto a los
// datos de su transacción contenedora. No es un agregado ni tiene repositorio de
// escritura: solo existe para consultas de lectura sobre la vista `gastos`.
type Gasto interface {
	filter.Filterable

	Operacion() string
	Cuota() string
	Fecha() time.Time
	Concepto() string
	Monto() quantity.Quantity
	Moneda() moneda.Moneda
	Total() quantity.Quantity
	Metodo() metodoperacion.MetodoDeOperacion
	Tasa() quantity.Quantity
	Registrado_por() string
	Registro() time.Time

	AsAProveedor() *GastoAProveedor
}

type GastoBase struct {
	operacion      string
	cuota          string
	fecha          time.Time
	concepto       string
	monto          quantity.Quantity
	moneda         moneda.Moneda
	metodo         metodoperacion.MetodoDeOperacion
	tasa           quantity.Quantity
	registrado_por string
	registro       time.Time
}

func (g GastoBase) Operacion() string        { return g.operacion }
func (g GastoBase) Cuota() string            { return g.cuota }
func (g GastoBase) Fecha() time.Time         { return g.fecha }
func (g GastoBase) Concepto() string         { return g.concepto }
func (g GastoBase) Monto() quantity.Quantity { return g.monto }
func (g GastoBase) Moneda() moneda.Moneda    { return g.moneda }
func (g GastoBase) Total() quantity.Quantity {
	if g.moneda == moneda.USD {
		return g.monto
	}
	return g.monto.HappyDiv(g.tasa)
}

func (g GastoBase) AsAProveedor() *GastoAProveedor           { return nil }
func (g GastoBase) Metodo() metodoperacion.MetodoDeOperacion { return g.metodo }
func (g GastoBase) Tasa() quantity.Quantity                  { return g.tasa }
func (g GastoBase) Registrado_por() string                   { return g.registrado_por }
func (g GastoBase) Registro() time.Time                      { return g.registro }
func (g GastoBase) FilterSpec() filter.Spec {
	return filter.Spec{
		"concepto": filter.TypeString,
		"cuota":    filter.TypeString,
	}
}

type GastoAProveedor struct {
	GastoBase
	Proveedor string
}

func (g *GastoAProveedor) AsAProveedor() *GastoAProveedor { return g }

type GastoFinder interface {
	Buscar(
		ctx context.Context,
		filter filter.Clause,
		paginator common.Paginator,
	) (*common.Paginated[Gasto], core.Error)
}
