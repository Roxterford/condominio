package gorm

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/Sanaruca/condominio/internal/core"
	gormAdapter "github.com/Sanaruca/condominio/internal/core/adapters/gorm"
	"github.com/Sanaruca/condominio/internal/core/common"
	"github.com/Sanaruca/condominio/internal/core/common/filter"
	"github.com/Sanaruca/condominio/internal/unidades/models/sujeto"
	"github.com/Sanaruca/condominio/internal/unidades/models/unidad"
	"github.com/Sanaruca/condominio/internal/unidades/models/unidad/estadounidad"
)

type GORMUnidadRepository struct {
	db *gorm.DB
	uf *unidad.UnidadFactory
	sf *sujeto.SujetoFactory
}

func NewGORMUnidadRepository(
	db *gorm.DB,
	unidadFactory *unidad.UnidadFactory,
	sujetoFactory *sujeto.SujetoFactory,
) unidad.UnidadRepository {

	if unidadFactory == nil {
		panic("unidad factory is nil")
	}

	if sujetoFactory == nil {
		panic("sujeto factory is nil")
	}

	return &GORMUnidadRepository{db, unidadFactory, sujetoFactory}
}

// ObtenerPorID implements [unidad.UnidadRepository].
func (r *GORMUnidadRepository) ObtenerPorID(ctx context.Context, id unidad.UnidadID) (*unidad.Unidad, core.Error) {

	unidad_row, err := gorm.G[UnidadInfo](r.db).
		Where("id = ?", id.String()).
		Preload("Contacto", nil).
		Preload("TitularPrimario", nil).
		Take(ctx)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if err != nil {
		return nil, core.WrapError(err)
	}

	unidad := unidad_row.ToDomainUnidad(r.uf, r.sf)

	return &unidad, nil
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
		Preload("Contacto", nil).
		Preload("TitularPrimario", nil).
		Order("codigo asc"). // TODO: ordenar para que no pase 1, 10, 2, 20 ...
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
		unidades = append(unidades, u.ToDomainUnidad(r.uf, r.sf))
	}

	return &common.Paginated[unidad.Unidad]{
		Data:  unidades,
		Total: int(total),
		Page:  p.Page,
		Pages: pages,
		Limit: p.Limit,
	}, nil
}
