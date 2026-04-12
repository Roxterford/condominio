package gorm

import (
	"context"
	"errors"

	"github.com/Sanaruca/condominio/internal/administracion/models/deuda"
	"github.com/Sanaruca/condominio/internal/administracion/types/estadodeuda"
	"github.com/Sanaruca/condominio/internal/core"
	gormAdapter "github.com/Sanaruca/condominio/internal/core/adapters/gorm"
	"github.com/Sanaruca/condominio/internal/core/common/filter"
	"gorm.io/gorm"
)

type GORMDeudaRepository struct {
	db           *gorm.DB
	deudaFactory *deuda.DeudaFactory
}

// Count implements [deuda.DeudaRepository].
func (r GORMDeudaRepository) Count(ctx context.Context, filter filter.Clause) (int, core.Error) {
	count, err := gorm.G[Deuda](r.db).Scopes(gormAdapter.GFilter(filter)).Count(ctx, "id")

	if err != nil {
		return 0, core.WrapError(err)
	}

	return int(count), nil
}

func NewGORMDeudaRepository(
	db *gorm.DB,
	deudaFactory *deuda.DeudaFactory,
) deuda.DeudaRepository {

	if deudaFactory == nil {
		panic("deudaFactory is nil")
	}

	return GORMDeudaRepository{db: db, deudaFactory: deudaFactory}
}

// GetLastDeudaWhereNotPagada implements [deuda.DeudaRepository].
func (r GORMDeudaRepository) GetLastDeudaWhereNotPagada(
	ctx context.Context,
	villa int,
) (*deuda.Deuda, core.Error) {
	_deuda, err := gorm.G[Deuda](r.db).Where(
		"villa = ? AND estado <> ?",
		villa,
		estadodeuda.Pagada,
	).Select("id", "deuda").Order("registro asc").Take(ctx)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if err != nil {
		return nil, core.WrapError(err)
	}

	destinos, err := gorm.G[DestinoDePago](r.db).Where("deuda = ?", _deuda.ID).Find(ctx)

	if err != nil {
		return nil, core.WrapError(err)
	}

	abonos := make([]deuda.Abono, len(destinos))
	for i, destino := range destinos {
		abonos[i] = *r.deudaFactory.AssembleAbono(destino.Pago, destino.Destinado, destino.Fecha)
	}

	return r.deudaFactory.Assemble(
		_deuda.ID,
		_deuda.Cuota,
		_deuda.Villa,
		_deuda.Monto,
		_deuda.Registro,
		abonos,
	)
}
