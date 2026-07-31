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
	"github.com/Sanaruca/condominio/internal/finanzas/models/transaccion"
)

type GORMTransaccionRepository struct {
	db *gorm.DB
	qf *quantity.QuantityFactory
	tf *transaccion.TransaccionFactory
}

func NewGORMTransaccionRepository(
	db *gorm.DB,
	quantityFactory *quantity.QuantityFactory,
	transaccionFactory *transaccion.TransaccionFactory,
) transaccion.TransaccionRepository {
	if quantityFactory == nil {
		panic("quantityFactory is nil")
	}
	if transaccionFactory == nil {
		panic("transaccionFactory is nil")
	}

	return &GORMTransaccionRepository{
		db: db,
		qf: quantityFactory,
		tf: transaccionFactory,
	}
}

func (r *GORMTransaccionRepository) Guardar(
	ctx context.Context,
	tx *transaccion.TransaccionFinanciera,
) core.Error {
	_, err := gorm.G[ITransaccion](r.db).Where("id = ?", tx.ID()).Select("id").Take(ctx)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return r.insertar(ctx, tx)
	} else if err != nil {
		return core.WrapError(err)
	}

	return r.actualizar(ctx, tx)
}

func (r *GORMTransaccionRepository) insertar(
	ctx context.Context,
	tx *transaccion.TransaccionFinanciera,
) core.Error {
	transaccionTable := mapToITransaccion(tx)
	movimientos := mapToIMovimientos(tx)

	err := r.db.Transaction(func(dbTX *gorm.DB) error {
		if err := gorm.G[ITransaccion](dbTX).Create(ctx, transaccionTable); err != nil {
			return err
		}

		if len(movimientos) > 0 {
			if err := gorm.G[IMovimiento](dbTX).CreateInBatches(ctx, &movimientos, 100); err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return core.WrapError(err)
	}

	return nil
}

func (r *GORMTransaccionRepository) actualizar(
	ctx context.Context,
	tx *transaccion.TransaccionFinanciera,
) core.Error {
	transaccionTable := mapToITransaccion(tx)

	err := r.db.Transaction(func(dbTX *gorm.DB) error {
		if _, err := gorm.G[ITransaccion](dbTX).Where("id = ?", tx.ID()).Updates(ctx, *transaccionTable); err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return core.WrapError(err)
	}

	return nil
}

func (r *GORMTransaccionRepository) ObtenerPorID(
	ctx context.Context,
	id string,
) (*transaccion.TransaccionFinanciera, core.Error) {
	dbTX, err := gorm.G[ITransaccion](r.db).Where("id = ?", id).First(ctx)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, transaccion.ErrTransaccionNoEncontrada
	}

	if err != nil {
		return nil, core.WrapError(err)
	}

	movs, err := gorm.G[IMovimiento](r.db).Where("transaccion_id = ?", id).Find(ctx)
	if err != nil {
		return nil, core.WrapError(err)
	}

	movimientos := make([]transaccion.Movimiento, len(movs))
	for i, m := range movs {
		movimientos[i] = toDomainMovimiento(m, r.qf)
	}

	return r.tf.Assemble(
		dbTX.ID,
		dbTX.Fecha,
		dbTX.Concepto,
		r.qf.Assemble(int64(dbTX.MontoTotal)),
		dbTX.Moneda,
		dbTX.Metodo,
		r.qf.Assemble(int64(dbTX.Tasa)),
		dbTX.RegistradoPor,
		dbTX.CuotaID,
		dbTX.Registro,
		movimientos,
	), nil
}

func (r *GORMTransaccionRepository) Obtener(
	ctx context.Context,
	filter filter.Clause,
	paginator common.Paginator,
) (*common.Paginated[transaccion.TransaccionFinanciera], core.Error) {
	paginator.Sanitize()

	rows, err := gorm.G[ITransaccion](r.db).
		Scopes(
			gormAdapter.GFilter(filter),
			gormAdapter.GPaginate(paginator),
		).
		Find(ctx)

	if err != nil {
		return nil, core.WrapError(err)
	}

	total, err := gorm.G[ITransaccion](r.db).
		Scopes(gormAdapter.GFilter(filter)).
		Count(ctx, "id")

	if err != nil {
		return nil, core.WrapError(err)
	}

	data := make([]transaccion.TransaccionFinanciera, len(rows))
	for i, t := range rows {
		movs, err := gorm.G[IMovimiento](r.db).Where("transaccion_id = ?", t.ID).Find(ctx)
		if err != nil {
			return nil, core.WrapError(err)
		}

		movimientos := make([]transaccion.Movimiento, len(movs))
		for j, m := range movs {
			movimientos[j] = toDomainMovimiento(m, r.qf)
		}

		data[i] = *r.tf.Assemble(
			t.ID,
			t.Fecha,
			t.Concepto,
			r.qf.Assemble(int64(t.MontoTotal)),
			t.Moneda,
			t.Metodo,
			r.qf.Assemble(int64(t.Tasa)),
			t.RegistradoPor,
			t.CuotaID,
			t.Registro,
			movimientos,
		)
	}

	return common.NewPaginated(data, int(total), paginator), nil
}

func (r *GORMTransaccionRepository) Count(
	ctx context.Context,
	filter filter.Clause,
) (int, core.Error) {
	count, err := gorm.G[ITransaccion](r.db).Scopes(gormAdapter.GFilter(filter)).Count(ctx, "id")
	if err != nil {
		return 0, core.WrapError(err)
	}

	return int(count), nil
}
