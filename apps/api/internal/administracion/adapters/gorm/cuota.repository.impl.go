package gorm

import (
	"context"
	"errors"

	"github.com/Sanaruca/condominio/internal/administracion/models/cuota"
	"github.com/Sanaruca/condominio/internal/core"
	gormAdapter "github.com/Sanaruca/condominio/internal/core/adapters/gorm"
	"github.com/Sanaruca/condominio/internal/core/common"
	"github.com/Sanaruca/condominio/internal/core/common/filter"
	"gorm.io/gorm"
)

type GORMCuotaRepository struct {
	db      *gorm.DB
	factory *cuota.CuotaFactory
}

func NewGORMCuotaRepository(db *gorm.DB, factory *cuota.CuotaFactory) cuota.CuotaRepository {
	if factory == nil {
		panic("factory is nil")
	}
	return &GORMCuotaRepository{db: db, factory: factory}
}

// ObtenerPorID implements [cuota.CuotaRepository].
func (r *GORMCuotaRepository) ObtenerPorID(
	ctx context.Context,
	id cuota.CuotaID,
) (cuota.Cuota, core.Error) {
	cuota, err := gorm.G[Cuota](
		r.db,
	).Preload("Proyecto", nil).
		Where("id = ?", id.String()).
		Take(ctx)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, core.WrapError(err)
	}

	return cuota.ToDomainCuota(r.factory), nil
}

func (r *GORMCuotaRepository) Obtener(
	ctx context.Context,
	filter filter.Clause,
	paginator common.Paginator,
) (*common.Paginated[cuota.Cuota], core.Error) {

	// Obtener registros paginados
	registros, err := gorm.G[Cuota](r.db).
		Scopes(gormAdapter.GFilter(filter), gormAdapter.GPaginate(paginator)).
		Preload("Proyecto", nil).
		Find(ctx)
	if err != nil {
		return nil, core.WrapError(err)
	}

	// Obtener total para la metadata de paginación
	total, err := gorm.G[Cuota](r.db).Scopes(gormAdapter.GFilter(filter)).Count(ctx, "id")
	if err != nil {
		return nil, core.WrapError(err)
	}

	// Cálculo de páginas optimizado
	pages := (int(total) + paginator.Limit - 1) / paginator.Limit

	// 5. Mapeo final a objetos de dominio
	cuotas := make([]cuota.Cuota, 0, len(registros))
	for _, c := range registros {
		// Se pasa el detalle (si no existe en el mapa, será el valor cero de la estructura)
		cuotas = append(cuotas, c.ToDomainCuota(r.factory))
	}

	return &common.Paginated[cuota.Cuota]{
		Data:  cuotas,
		Total: int(total),
		Page:  paginator.Page,
		Pages: pages,
		Limit: paginator.Limit,
	}, nil
}
