package gorm

import (
	"context"
	"errors"

	gormPkg "gorm.io/gorm"

	"github.com/Sanaruca/condominio/internal/core"
	"github.com/Sanaruca/condominio/internal/core/common/quantity"
	"github.com/Sanaruca/condominio/internal/unidades/models/unidad"
)

type GORMUnidadEstadisticasFinder struct {
	db *gormPkg.DB
	qf *quantity.QuantityFactory
}

func NewGORMUnidadEstadisticasFinder(db *gormPkg.DB, qf *quantity.QuantityFactory) unidad.EstadisticasFinder {
	if qf == nil {
		panic("qf is nil")
	}
	return &GORMUnidadEstadisticasFinder{db, qf}
}

func (f *GORMUnidadEstadisticasFinder) Obtener(ctx context.Context) (*unidad.Estadisticas, core.Error) {
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
	err := f.db.WithContext(ctx).Table("unidades_totales").Take(&result).Error

	if errors.Is(err, gormPkg.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, core.WrapError(err)
	}

	return &unidad.Estadisticas{
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
