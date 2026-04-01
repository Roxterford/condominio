package command

import (
	"github.com/Sanaruca/condominio/internal/administracion/models/proveedor"
	"github.com/Sanaruca/condominio/internal/core"
	"github.com/Sanaruca/condominio/internal/core/adapters/ozzo"
	"github.com/Sanaruca/condominio/internal/core/common"
	"github.com/Sanaruca/condominio/internal/core/context"
	"github.com/Sanaruca/condominio/internal/core/usecase"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type RegistrarProveedorDTO struct {
	Rif       string
	Nombre    string
	Email     string
	Telefono  string
	Direccion *string
}

type RegistrarProveedor usecase.Handler[context.AdminContext, RegistrarProveedorDTO, *proveedor.Proveedor]

type registrarProveedor struct {
	repo         proveedor.ProveedorRepository
	factory      proveedor.ProveedorFactory
	emailFactory common.EmailFactory
	phoneFactory common.PhoneFactory
}

func NewRegistrarProveedor(repo proveedor.ProveedorRepository) RegistrarProveedor {
	if repo == nil {
		panic("repo is nil")
	}
	return &registrarProveedor{repo: repo}
}

func (uc *registrarProveedor) Exec(
	ctx context.AdminContext,
	input RegistrarProveedorDTO,
) (*proveedor.Proveedor, core.Error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	existe, err := uc.repo.ExistePorRif(ctx, input.Rif)
	if err != nil {
		return nil, err
	}
	if existe {
		return nil, proveedor.ErrorProveedorDuplicado
	}

	proveedor, err := uc.factory.Nuevo(
		input.Rif,
		input.Nombre,
		input.Email,
		input.Telefono,
		input.Direccion,
	)

	if err != nil {
		return nil, err
	}

	if err := uc.repo.Guardar(ctx, proveedor); err != nil {
		return nil, err
	}

	return proveedor, nil
}

func (dto *RegistrarProveedorDTO) Validate() core.Error {
	err := validation.ValidateStruct(
		dto,
		validation.Field(&dto.Rif, validation.Required, validation.Length(1, 20)),
		validation.Field(&dto.Nombre, validation.Required, validation.Length(1, 100)),
		validation.Field(&dto.Email, validation.Required, is.Email),
		validation.Field(&dto.Telefono, validation.Required, validation.Length(1, 20)),
		validation.Field(&dto.Direccion, validation.NilOrNotEmpty, validation.Length(1, 200)),
	)

	if err != nil {
		return ozzo.FirstOzzoErrorAdapter(dto, err)
	}

	return nil
}
