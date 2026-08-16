package gorm

import (
	"context"

	"gorm.io/gorm"

	"github.com/Sanaruca/condominio/internal/core"
	gormAdapter "github.com/Sanaruca/condominio/internal/core/adapters/gorm"
	"github.com/Sanaruca/condominio/internal/core/common"
	"github.com/Sanaruca/condominio/internal/core/common/filter"
	"github.com/Sanaruca/condominio/internal/core/common/quantity"
	"github.com/Sanaruca/condominio/internal/finanzas/models/operacion"
)

type GORMGastoFinder struct {
	db *gorm.DB
	qf *quantity.QuantityFactory
}

func NewGORMGastoFinder(
	db *gorm.DB,
	quantityFactory *quantity.QuantityFactory,
) operacion.GastoFinder {
	if db == nil {
		panic("db is nil")
	}
	if quantityFactory == nil {
		panic("quantityFactory is nil")
	}

	return &GORMGastoFinder{
		db: db,
		qf: quantityFactory,
	}
}

// Buscar implements [operacion.GastoFinder].
func (g *GORMGastoFinder) Buscar(
	ctx context.Context,
	filter filter.Clause,
	paginator common.Paginator,
) (*common.Paginated[operacion.Gasto], core.Error) {
	paginator.Sanitize()

	rows, err := gorm.G[Gasto](g.db).Scopes(
		gormAdapter.GFilter(filter),
		gormAdapter.GPaginate(paginator),
	).Find(ctx)
	if err != nil {
		return nil, core.WrapError(err)
	}

	total, err := gorm.G[Gasto](g.db).Scopes(gormAdapter.GFilter(filter)).Count(ctx, "operacion")
	if err != nil {
		return nil, core.WrapError(err)
	}

	data := make([]operacion.Gasto, len(rows))
	for i, row := range rows {
		data[i] = mapToGasto(row, g.qf)
	}

	return common.NewPaginated(data, int(total), paginator), nil
}
