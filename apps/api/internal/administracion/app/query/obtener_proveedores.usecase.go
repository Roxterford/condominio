package query

import (
	"github.com/Sanaruca/condominio/internal/administracion/models/proveedor"
	"github.com/Sanaruca/condominio/internal/core"
	"github.com/Sanaruca/condominio/internal/core/context"
	"github.com/Sanaruca/condominio/internal/core/usecase"
)

// TODO: Cambiar a AdminContext cuando la autenticación funcione
type ObtenerProveedores usecase.Handler[context.BaseContext, any, []*proveedor.Proveedor]

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
	_ any,
) ([]*proveedor.Proveedor, core.Error) {
	proveedores, err := uc.repo.ObtenerTodos(ctx)
	if err != nil {
		return nil, err
	}
	return proveedores, nil
}
