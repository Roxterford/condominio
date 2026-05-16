package gorm

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/Sanaruca/condominio/internal/administracion/models/gasto"
	"github.com/Sanaruca/condominio/internal/core"
	gormAdapter "github.com/Sanaruca/condominio/internal/core/adapters/gorm"
	"github.com/Sanaruca/condominio/internal/core/common"
	"github.com/Sanaruca/condominio/internal/core/common/quantity"
)

type GORMGastoRepository struct {
	db      *gorm.DB
	factory *gasto.GastoFactory
	qf      *quantity.QuantityFactory
}

func NewGORMGastoRepository(
	db *gorm.DB,
	gastoFactory *gasto.GastoFactory,
	quantityFactory *quantity.QuantityFactory,
) gasto.GastoRepository {

	if gastoFactory == nil {
		panic("gastoFactory is nil")
	}

	if quantityFactory == nil {
		panic("quantityFactory is nil")
	}

	return &GORMGastoRepository{
		db:      db,
		factory: gastoFactory,
		qf:      quantityFactory,
	}
}

func (r *GORMGastoRepository) ObtenerTodos(
	ctx context.Context,
	paginator common.Paginator,
) (*common.Paginated[gasto.Gasto], core.Error) {

	total, err := gorm.G[Gasto](
		r.db,
	).
		Count(ctx, "*")

	if err != nil {
		return nil, core.WrapError(err)
	}

	db_gastos, err := gorm.G[Gasto](r.db).
		Scopes(gormAdapter.GPaginate(paginator)).
		Find(ctx)

	if err != nil {
		return nil, core.WrapError(err)
	}

	gastos := make([]gasto.Gasto, len(db_gastos))
	for i, db_gasto := range db_gastos {
		gastos[i] = *db_gasto.ToDomainGasto(r.factory, r.qf)
	}

	return common.NewPaginated(gastos, int(total), paginator), nil
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

func (r *GORMGastoRepository) ObtenerPorIDs(
	ctx context.Context,
	ids []gasto.GastoID,
) ([]gasto.Gasto, core.Error) {
	if len(ids) == 0 {
		return []gasto.Gasto{}, nil
	}

	idStrings := make([]string, len(ids))
	for i, id := range ids {
		idStrings[i] = string(id)
	}

	var models []GastoView
	if err := r.db.WithContext(ctx).Where("id IN ?", idStrings).Find(&models).Error; err != nil {
		return nil, core.WrapError(err)
	}

	gastos := make([]gasto.Gasto, len(models))
	for i, model := range models {
		gastos[i] = *r.toGasto(&model)
	}

	return gastos, nil
}

func (r *GORMGastoRepository) Actualizar(
	ctx context.Context,
	gastoEntity gasto.Gasto,
) core.Error {
	model := toGastoTable(&gastoEntity)
	if err := r.db.WithContext(ctx).Model(model).Updates(map[string]any{
		"cuota": model.Cuota,
	}).Error; err != nil {
		return core.WrapError(err)
	}
	return nil
}

func toGastoTable(g *gasto.Gasto) *Gasto {
	return &Gasto{
		ID:             string(g.ID()),
		Proveedor:      g.Proveedor(),
		Cuota:          g.Cuota(),
		Monto:          int(g.Monto().Value()),
		Moneda:         g.Moneda(),
		Tasa:           int(g.Tasa().Value()),
		Fecha:          g.Fecha(),
		Descripcion:    g.Descripcion(),
		Registrado_por: g.Audit().CreatedBy,
		Registro:       g.Audit().CreatedAt,
	}
}

func (r *GORMGastoRepository) toGasto(table *GastoView) *gasto.Gasto {
	return r.factory.Assemble(
		table.ID,
		table.Concepto,
		table.Proveedor,
		table.Cuota,
		r.qf.Assemble(int64(table.Monto)),
		table.Moneda,
		r.qf.Assemble(int64(table.Tasa)),
		table.Fecha,
		table.Descripcion,
		table.Registrado_por,
		table.Registro,
	)
}
