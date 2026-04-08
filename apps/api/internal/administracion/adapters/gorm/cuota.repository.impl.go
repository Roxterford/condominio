package gorm

import (
	"context"

	"github.com/Sanaruca/condominio/internal/administracion/models/cuota"
	"github.com/Sanaruca/condominio/internal/administracion/types/tipodecuota"
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

func (r *GORMCuotaRepository) Obtener(
	ctx context.Context,
	filter filter.Clause,
	paginator common.Paginator,
) (*common.Paginated[cuota.Cuota], core.Error) {

	// Obtener registros paginados
	registros, err := gorm.G[Cuota](r.db).
		Scopes(gormAdapter.GFilter(filter), gormAdapter.GPaginate(paginator)).
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

	// Filtrar IDs para la carga de Proyectos (Eager Loading manual)
	ce_refs := make([]string, 0, len(registros))
	for _, c := range registros {
		if c.Tipo == tipodecuota.Especial {
			ce_refs = append(ce_refs, c.ID)
		}
	}

	detalles := make(map[string]cuota.Proyecto)

	// Optimización: Solo consultar si hay referencias
	if len(ce_refs) > 0 {
		proyectos, err := gorm.G[Proyecto](r.db).Where("cuota IN ?", ce_refs).Find(ctx)
		if err != nil {
			return nil, core.WrapError(err)
		}

		// Pre-asignar capacidad al mapa
		detalles = make(map[string]cuota.Proyecto, len(proyectos))
		proyectoFactory := r.factory.ProyectoFactory()
		for _, p := range proyectos {
			detalles[p.Cuota] = p.ToDomainProyecto(proyectoFactory)
		}
	}

	// 5. Mapeo final a objetos de dominio
	cuotas := make([]cuota.Cuota, 0, len(registros))
	for _, c := range registros {
		// Se pasa el detalle (si no existe en el mapa, será el valor cero de la estructura)
		cuotas = append(cuotas, c.ToDomainCuota(r.factory, detalles[c.ID]))
	}

	return &common.Paginated[cuota.Cuota]{
		Data:  cuotas,
		Total: int(total),
		Page:  paginator.Page,
		Pages: pages,
		Limit: paginator.Limit,
	}, nil
}
