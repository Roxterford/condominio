package gorm

import (
	"context"
	"errors"

	"github.com/Sanaruca/condominio/internal/core"
	gormAdapter "github.com/Sanaruca/condominio/internal/core/adapters/gorm"
	"github.com/Sanaruca/condominio/internal/core/common"
	"github.com/Sanaruca/condominio/internal/core/common/filter"
	"github.com/Sanaruca/condominio/internal/unidades/models/unidad"
	"github.com/Sanaruca/condominio/internal/unidades/models/unidad/estadounidad"
	"gorm.io/gorm"
)

type GORMUnidadRepository struct {
	db      *gorm.DB
	factory *unidad.UnidadFactory
}

func NewGORMUnidadRepository(db *gorm.DB, factory *unidad.UnidadFactory) unidad.UnidadRepository {

	if factory == nil {
		panic("factory is nil")
	}

	return &GORMUnidadRepository{db, factory}
}

func (r *GORMUnidadRepository) Exists(
	ctx context.Context,
	unidadID unidad.UnidadID,
) (bool, core.Error) {

	var u Unidad
	err := r.db.WithContext(ctx).Where("codigo = ?", unidadID.String()).Select("id").Take(&u).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}

	if err != nil {
		return false, core.WrapError(err)
	}

	return true, nil
}

func (r *GORMUnidadRepository) ObtenerTodas(ctx context.Context) ([]string, core.Error) {
	var unidades []Unidad
	if err := r.db.WithContext(ctx).Select("codigo").Find(&unidades).Error; err != nil {
		return nil, core.WrapError(err)
	}

	codigos := make([]string, len(unidades))
	for i, u := range unidades {
		codigos[i] = u.Codigo
	}

	return codigos, nil
}

func (r *GORMUnidadRepository) ObtenerEstado(
	ctx context.Context,
	unidadCodigo string,
) (estadounidad.EstadoDeUnidad, core.Error) {
	var u Unidad
	if err := r.db.WithContext(ctx).Where("codigo = ?", unidadCodigo).Take(&u).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", unidad.ErrUnidadNoEncontrada
		}
		return "", core.WrapError(err)
	}

	return estadounidad.EstadoDeUnidad(u.Estado), nil
}

func (r *GORMUnidadRepository) Obtener(
	ctx context.Context,
	f filter.Clause,
	p common.Paginator,
) (*common.Paginated[unidad.Unidad], core.Error) {
	p.Sanitize()
	// Obtener registros paginados
	registros, err := gorm.G[UnidadInfo](r.db).
		Scopes(gormAdapter.GFilter(f), gormAdapter.GPaginate(p)).
		Find(ctx)
	if err != nil {
		return nil, core.WrapError(err)
	}

	// Obtener total para la metadata de paginación
	total, err := gorm.G[Unidad](r.db).Scopes(gormAdapter.GFilter(f)).Count(ctx, "id")
	if err != nil {
		return nil, core.WrapError(err)
	}

	// Cálculo de páginas optimizado
	pages := (int(total) + p.Limit - 1) / p.Limit

	// 5. Mapeo final a objetos de dominio
	unidades := make([]unidad.Unidad, 0, len(registros))
	for _, u := range registros {
		// Se pasa el detalle (si no existe en el mapa, será el valor cero de la estructura)
		unidades = append(unidades, u.ToDomainUnidad(r.factory))
	}

	return &common.Paginated[unidad.Unidad]{
		Data:  unidades,
		Total: int(total),
		Page:  p.Page,
		Pages: pages,
		Limit: p.Limit,
	}, nil
}
