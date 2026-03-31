package gorm

import (
	"context"
	"errors"

	"github.com/Sanaruca/condominio/internal/administracion/models/gasto"
	"github.com/Sanaruca/condominio/internal/core"
	"gorm.io/gorm"
)

type GORMGastoRepository struct {
	db      *gorm.DB
	factory *gasto.GastoFactory
}

func NewGORMGastoRepository(
	db *gorm.DB,
	factory *gasto.GastoFactory,
) gasto.GastoRepository {
	return &GORMGastoRepository{db: db, factory: factory}
}

func (r *GORMGastoRepository) Guardar(
	ctx context.Context,
	gastoEntity gasto.Gasto,
) (*gasto.GastoID, core.Error) {
	model := toGastoTable(&gastoEntity)
	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return nil, core.WrapError(err)
	}
	id := gasto.GastoID(model.ID)
	return &id, nil
}

func (r *GORMGastoRepository) ObtenerPorID(
	ctx context.Context,
	id gasto.GastoID,
) (*gasto.Gasto, core.Error) {
	var model GastoView
	if err := r.db.WithContext(ctx).Where("id = ?", string(id)).Take(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gasto.ErrGastoNotFound
		}
		return nil, core.WrapError(err)
	}

	return r.toGasto(&model), nil
}

func toGastoTable(g *gasto.Gasto) *Gasto {
	return &Gasto{
		ID:             string(g.ID()),
		Proveedor:      g.Proveedor(),
		Cuota:          g.Cuota(),
		Monto:          g.Monto(),
		Moneda:         g.Moneda(),
		Tasa:           g.Tasa(),
		Fecha:          g.Fecha(),
		Descripcion:    g.Descripcion(),
		Registrado_por: g.Audit().CreatedBy,
		Registro:       g.Audit().CreatedAt,
	}
}

func (r *GORMGastoRepository) toGasto(table *GastoView) *gasto.Gasto {
	return r.factory.Assemble(
		table.ID,
		table.Proveedor,
		table.Cuota,
		table.Monto,
		table.Moneda,
		table.Tasa,
		table.Fecha,
		table.Descripcion,
		table.Registrado_por,
		table.Registro,
	)
}
