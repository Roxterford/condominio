package query

import (
	"github.com/Sanaruca/condominio/internal/administracion/models/proveedor"
	"github.com/Sanaruca/condominio/internal/core"
	"github.com/Sanaruca/condominio/internal/core/context"
	"github.com/Sanaruca/condominio/internal/core/exception"
	"github.com/Sanaruca/condominio/internal/core/usecase"
)

type ObtenerProveedorDTO struct {
	ID string
}

type ObtenerProveedor usecase.Handler[context.AdminContext, ObtenerProveedorDTO, *proveedor.Proveedor]

type obtenerProveedor struct {
	repo proveedor.ProveedorRepository
}

func NewObtenerProveedor(repo proveedor.ProveedorRepository) ObtenerProveedor {
	if repo == nil {
		panic("repo is nil")
	}
	return &obtenerProveedor{repo: repo}
}

func (uc *obtenerProveedor) Exec(
	ctx context.AdminContext,
	input ObtenerProveedorDTO,
) (*proveedor.Proveedor, core.Error) {
	if input.ID == "" {
		return nil, exception.New(exception.INVALID_ARGUMENT, "ID es requerido")
	}

	proveedor, err := uc.repo.ObtenerPorID(ctx, input.ID)
	if err != nil {
		return nil, err
	}

	return proveedor, nil
}
