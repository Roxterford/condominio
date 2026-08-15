package gorm

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/Sanaruca/condominio/internal/core"
	gormAdapter "github.com/Sanaruca/condominio/internal/core/adapters/gorm"
	"github.com/Sanaruca/condominio/internal/core/common"
	"github.com/Sanaruca/condominio/internal/core/common/filter"
	"github.com/Sanaruca/condominio/internal/core/common/quantity"
	"github.com/Sanaruca/condominio/internal/finanzas/models/operacion"
)

type GORMOperacionRepository struct {
	db *gorm.DB
	qf *quantity.QuantityFactory
	f  *operacion.OperacionFactory
}

func NewGORMOperacionRepository(
	db *gorm.DB,
	quantityFactory *quantity.QuantityFactory,
	operacionFactory *operacion.OperacionFactory,
) operacion.OperacionRepository {
	if db == nil {
		panic("db is nil")
	}
	if quantityFactory == nil {
		panic("quantityFactory is nil")
	}
	if operacionFactory == nil {
		panic("operacionFactory is nil")
	}

	return &GORMOperacionRepository{
		db: db,
		qf: quantityFactory,
		f:  operacionFactory,
	}
}

// Guardar persiste una operacion. Las operaciones son inmutables: si el id ya
// existe, la escritura es un no-op.
func (r *GORMOperacionRepository) Guardar(
	ctx context.Context,
	op *operacion.Operacion,
) core.Error {
	_, err := gorm.G[Operacion](r.db).Where("id = ?", op.ID()).Select("id").Take(ctx)

	if err == nil {
		return nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return core.WrapError(err)
	}

	if err := gorm.G[Operacion](r.db).Create(ctx, mapToOperacion(op)); err != nil {
		return core.WrapError(err)
	}

	return nil
}

func (r *GORMOperacionRepository) ObtenerPorID(
	ctx context.Context,
	id string,
) (*operacion.Operacion, core.Error) {
	row, err := gorm.G[Operacion](r.db).Where("id = ?", id).First(ctx)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, operacion.ErrOperacionNoEncontrada
	}
	if err != nil {
		return nil, core.WrapError(err)
	}

	return toDomainOperacion(row, r.qf, r.f), nil
}

func (r *GORMOperacionRepository) Obtener(
	ctx context.Context,
	filter filter.Clause,
	paginator common.Paginator,
) (*common.Paginated[operacion.Operacion], core.Error) {
	paginator.Sanitize()

	rowsQuery := gorm.G[Operacion](r.db).Scopes(
		gormAdapter.GFilter(filter),
		gormAdapter.GPaginate(paginator),
	)

	rows, err := rowsQuery.Find(ctx)
	if err != nil {
		return nil, core.WrapError(err)
	}

	countQuery := gorm.G[Operacion](r.db).Scopes(gormAdapter.GFilter(filter))

	total, err := countQuery.Count(ctx, "id")
	if err != nil {
		return nil, core.WrapError(err)
	}

	data := make([]operacion.Operacion, len(rows))
	for i, row := range rows {
		data[i] = *toDomainOperacion(row, r.qf, r.f)
	}

	return common.NewPaginated(data, int(total), paginator), nil
}

func (r *GORMOperacionRepository) Count(
	ctx context.Context,
	filter filter.Clause,
) (int, core.Error) {
	count, err := gorm.G[Operacion](r.db).Scopes(gormAdapter.GFilter(filter)).Count(ctx, "id")
	if err != nil {
		return 0, core.WrapError(err)
	}

	return int(count), nil
}
