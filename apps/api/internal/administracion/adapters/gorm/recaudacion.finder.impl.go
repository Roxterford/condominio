package gorm

import (
	"context"
	"errors"

	"github.com/Sanaruca/condominio/internal/administracion/models/cuota"
	"github.com/Sanaruca/condominio/internal/core"
	"github.com/Sanaruca/condominio/internal/core/common"
	"github.com/Sanaruca/condominio/internal/core/common/filter"
	"github.com/Sanaruca/condominio/internal/core/common/quantity"
	"gorm.io/gorm"
)

type GROMRecaudacionFinder struct {
	db *gorm.DB
	qf *quantity.QuantityFactory
}

func NewGROMRecaudacionFinder(db *gorm.DB, qf *quantity.QuantityFactory) cuota.RecaudacionFinder {

	if qf == nil {
		panic("qf is nil")
	}

	return &GROMRecaudacionFinder{db, qf}
}

// Buscar implements [cuota.RecaudacionFinder].
func (f *GROMRecaudacionFinder) Buscar(
	ctx context.Context,
	filter filter.Clause,
) (*common.Paginated[cuota.Recaudacion], core.Error) {
	panic("unimplemented")
}

// ObtenerRecaudacion implements [cuota.RecaudacionFinder].
func (f *GROMRecaudacionFinder) ObtenerRecaudacion(
	ctx context.Context,
	id cuota.CuotaID,
) (*cuota.Recaudacion, core.Error) {
	recaudacion, err := gorm.G[Recaudacion](f.db).Where("cuota = ?", id).Take(ctx)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if err != nil {
		return nil, core.WrapError(err)
	}

	return recaudacion.ToDomainRecaudacion(f.qf), nil
}
