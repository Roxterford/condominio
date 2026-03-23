package command

import (
	"github.com/Sanaruca/condominio/internal/administracion/models/proveedor"
	"github.com/Sanaruca/condominio/internal/core"
	"github.com/Sanaruca/condominio/internal/core/context"
	"github.com/Sanaruca/condominio/internal/core/usecase"
)

type EliminarProveedorDTO struct {
	ID string
}

type EliminarProveedor usecase.WithContextInput[context.AdminContext, EliminarProveedorDTO]

type eliminarProveedor struct {
	repo proveedor.ProveedorRepository
}

func NewEliminarProveedor(repo proveedor.ProveedorRepository) EliminarProveedor {
	if repo == nil {
		panic("repo is nil")
	}
	return &eliminarProveedor{repo: repo}
}

func (uc *eliminarProveedor) Exec(
	ctx context.AdminContext,
	input EliminarProveedorDTO,
) (any, core.Error) {
	if input.ID == "" {
		return nil, core.NewInvalidArgumentError("ID es requerido")
	}

	existe, err := uc.repo.ExistePorID(ctx, input.ID)
	if err != nil {
		return nil, err
	}
	if !existe {
		return nil, proveedor.ErrProveedorNoEncontrado
	}

	// TODO: Verificar si el proveedor tiene contratos activos antes de eliminar

	if err := uc.repo.Eliminar(ctx, input.ID); err != nil {
		return nil, err
	}

	return nil, nil
}
