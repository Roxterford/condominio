package gorm

import (
	"context"
	"errors"

	"github.com/Sanaruca/condominio/internal/core"
	"github.com/Sanaruca/condominio/internal/core/common"
	"github.com/Sanaruca/condominio/internal/core/common/filter"
	unidadPkg "github.com/Sanaruca/condominio/internal/unidades/models/unidad"
	"github.com/Sanaruca/condominio/internal/unidades/models/unidad/estadounidad"
	gormPkg "gorm.io/gorm"
)

type GORMUnidadRepository struct {
	db *gormPkg.DB
}

func NewGORMUnidadRepository(db *gormPkg.DB) unidadPkg.UnidadRepository {
	return &GORMUnidadRepository{
		db: db,
	}
}

func (r *GORMUnidadRepository) Exists(
	ctx context.Context,
	unidadID unidadPkg.UnidadID,
) (bool, core.Error) {

	var u Unidad
	err := r.db.WithContext(ctx).Where("codigo = ?", unidadID.String()).Select("id").Take(&u).Error

	if errors.Is(err, gormPkg.ErrRecordNotFound) {
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
		if errors.Is(err, gormPkg.ErrRecordNotFound) {
			return "", unidadPkg.ErrUnidadNoEncontrada
		}
		return "", core.WrapError(err)
	}

	return estadounidad.EstadoDeUnidad(u.Estado), nil
}

func (r *GORMUnidadRepository) Obtener(
	ctx context.Context,
	f filter.Clause,
	p common.Paginator,
) (*common.Paginated[unidadPkg.Unidad], core.Error) {
	return nil, nil
}
