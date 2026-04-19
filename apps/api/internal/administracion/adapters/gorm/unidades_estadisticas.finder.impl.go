package gorm

import (
	"context"
	"errors"

	"github.com/Sanaruca/condominio/internal/administracion/models/unidades"
	"github.com/Sanaruca/condominio/internal/core"
	"github.com/Sanaruca/condominio/internal/core/common/quantity"
	gormPkg "gorm.io/gorm"
)

type GORMUnidadesEstadisticasFinder struct {
	db *gormPkg.DB
	qf *quantity.QuantityFactory
}

func NewGORMUnidadesEstadisticasFinder(db *gormPkg.DB, qf *quantity.QuantityFactory) unidades.UnidadesEstadisticasFinder {
	if qf == nil {
		panic("qf is nil")
	}
	return &GORMUnidadesEstadisticasFinder{db, qf}
}

func (f *GORMUnidadesEstadisticasFinder) Obtener(ctx context.Context) (*unidades.UnidadesEstadisticas, core.Error) {
	var result struct {
		TotalUnidades         int
		UnidadesActivas       int
		UnidadesInhabitadas   int
		UnidadesExentas       int
		UnidadesEnLitigio     int
		UnidadesSuspendidas   int
		UnidadesPreventa      int
		UnidadesConPendientes int
		UnidadesSolventes     int
		TotalPendiente        int
		TotalAsignado         int
	}
	err := f.db.WithContext(ctx).Table("unidades_estadisticas").Take(&result).Error

	if errors.Is(err, gormPkg.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, core.WrapError(err)
	}

	return &unidades.UnidadesEstadisticas{
		TotalUnidades:         result.TotalUnidades,
		UnidadesActivas:       result.UnidadesActivas,
		UnidadesInhabitadas:   result.UnidadesInhabitadas,
		UnidadesExentas:       result.UnidadesExentas,
		UnidadesEnLitigio:     result.UnidadesEnLitigio,
		UnidadesSuspendidas:   result.UnidadesSuspendidas,
		UnidadesPreventa:      result.UnidadesPreventa,
		UnidadesConPendientes: result.UnidadesConPendientes,
		UnidadesSolventes:     result.UnidadesSolventes,
		TotalPendiente:        f.qf.Assemble(int64(result.TotalPendiente)),
		TotalAsignado:         f.qf.Assemble(int64(result.TotalAsignado)),
	}, nil
}
