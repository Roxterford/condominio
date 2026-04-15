package gorm

import (
	"context"
	"errors"

	"github.com/Sanaruca/condominio/internal/administracion/models/villa"
	"github.com/Sanaruca/condominio/internal/core"
	"github.com/Sanaruca/condominio/internal/core/common/quantity"
	"gorm.io/gorm"
)

type GORMVillasTotalesFinder struct {
	db *gorm.DB
	qf *quantity.QuantityFactory
}

func NewGORMVillasTotalesFinder(db *gorm.DB, qf *quantity.QuantityFactory) villa.VillasTotalesFinder {

	if qf == nil {
		panic("qf is nil")
	}

	return &GORMVillasTotalesFinder{db, qf}
}

// ObtenerVillasTotales implements [villa.VillasTotalesFinder].
func (f *GORMVillasTotalesFinder) ObtenerVillasTotales(
	ctx context.Context,
) (*villa.VillasTotales, core.Error) {
	result, err := gorm.G[VillasTotales](f.db).Take(ctx)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if err != nil {
		return nil, core.WrapError(err)
	}

	return f.toDomain(result), nil
}

func (f *GORMVillasTotalesFinder) toDomain(t VillasTotales) *villa.VillasTotales {
	return &villa.VillasTotales{
		TotalVillas:         t.TotalVillas,
		VillasActivas:       t.VillasActivas,
		VillasInhabitadas:   t.VillasInhabitadas,
		VillasExentas:       t.VillasExentas,
		VillasEnLitigio:     t.VillasEnLitigio,
		VillasSuspendidas:   t.VillasSuspendidas,
		VillasPreventa:      t.VillasPreventa,
		VillasConPendientes: t.VillasConPendientes,
		VillasSolventes:     t.VillasSolventes,
		TotalPendiente:      f.qf.Assemble(int64(t.TotalPendiente)),
		TotalAsignado:       f.qf.Assemble(int64(t.TotalAsignado)),
	}
}
