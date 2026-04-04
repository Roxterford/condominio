package query

import (
	"github.com/Sanaruca/condominio/internal/administracion/models/proveedor"
	"github.com/Sanaruca/condominio/internal/core"
	"github.com/Sanaruca/condominio/internal/core/common/filter"
	"github.com/Sanaruca/condominio/internal/core/context"
	"github.com/Sanaruca/condominio/internal/core/usecase"
)

type ObtenerProveedoresDTO struct {
	Filter *filter.Filter[proveedor.Proveedor]
}

// TODO: Cambiar a AdminContext cuando la autenticación funcione
type ObtenerProveedores usecase.Handler[context.BaseContext, *ObtenerProveedoresDTO, []proveedor.Proveedor]

type obtenerProveedores struct {
	repo proveedor.ProveedorRepository
}

func NewObtenerProveedores(repo proveedor.ProveedorRepository) ObtenerProveedores {
	if repo == nil {
		panic("repo is nil")
	}
	return &obtenerProveedores{repo: repo}
}

func (uc *obtenerProveedores) Exec(
	ctx context.BaseContext,
	input *ObtenerProveedoresDTO,
) ([]proveedor.Proveedor, core.Error) {

	var clause filter.Clause

	if input != nil && input.Filter != nil {
		var err error
		clause, err = input.Filter.Build()
		if err != nil {
			return nil, core.WrapError(err)
		}
	}

	proveedores, err := uc.repo.Obtener(ctx, clause)
	if err != nil {
		return nil, err
	}
	return proveedores, nil
}
