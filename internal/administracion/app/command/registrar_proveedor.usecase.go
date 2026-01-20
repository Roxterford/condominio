package command

import (
	"errors"
	"regexp"
	"strings"
	"time"

	"github.com/Sanaruca/condominio/internal/administracion"
	"github.com/Sanaruca/condominio/internal/administracion/app/query"
	"github.com/Sanaruca/condominio/internal/administracion/types/tipodeproveedor"
	"github.com/Sanaruca/condominio/internal/core"
	"github.com/Sanaruca/condominio/internal/core/context"
	"github.com/Sanaruca/condominio/internal/core/usecase"
	"github.com/go-ozzo/ozzo-validation/is"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/lucsky/cuid"
	"gorm.io/gorm"
)

type RegistrarProveedorDTO struct {
	Rif       string
	Nombre    string
	Tipo      tipodeproveedor.TipoDeProveedor
	Email     string
	Telefono  string
	Direccion *string
}

type RegistrarProveedor usecase.Handler[context.AdminContext, RegistrarProveedorDTO, *administracion.Proveedor]

func NewRegistrarProveedor() RegistrarProveedor {
	return &registrarProveedor{}
}

type registrarProveedor struct {
	obtenerProveedor query.ObtenerProveedor
}

func (uc *registrarProveedor) Exec(
	ctx context.AdminContext,
	dto RegistrarProveedorDTO,
) (*administracion.Proveedor, core.Error) {

	if err := dto.Validate(); err != nil {
		return nil, err
	}

	proveedor := administracion.Proveedor{
		ID:            cuid.New(),
		Rif:           dto.Rif,
		Nombre:        dto.Nombre,
		Tipo:          dto.Tipo,
		Email:         dto.Email,
		Telefono:      dto.Telefono,
		Direccion:     dto.Direccion,
		Registro:      time.Now(),
		Actualizacion: time.Now(),
	}

	err := gorm.G[administracion.Proveedor](ctx.DB).Create(ctx.Context(), &proveedor)

	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return nil, administracion.ErrorProveedorDuplicado
	}

	if err != nil {
		return nil, core.WrapError(err)
	}

	return uc.obtenerProveedor.Exec(ctx.BaseContext, query.ObtenerProveedorDTO{
		ProveedorID: proveedor.ID,
	})
}

func (dto *RegistrarProveedorDTO) Validate() core.Error {

	dto.Rif = strings.TrimSpace(dto.Rif)
	dto.Nombre = strings.TrimSpace(dto.Nombre)
	dto.Tipo = tipodeproveedor.TipoDeProveedor(strings.TrimSpace(string(dto.Tipo)))
	dto.Email = strings.TrimSpace(dto.Email)
	dto.Telefono = strings.TrimSpace(dto.Telefono)
	if dto.Direccion != nil {
		*dto.Direccion = strings.TrimSpace(*dto.Direccion)
	}

	// RIF
	if dto.Rif == "" {
		return core.NewValidationError("El RIF es requerido")
	}
	if len(dto.Rif) > 15 {
		return core.NewValidationError("El RIF no puede tener mas de 15 caracteres")
	}

	// Nombre
	if dto.Nombre == "" {
		return core.NewValidationError("El nombre es requerido")
	}

	// Tipo
	if dto.Tipo == "" {
		return core.NewValidationError("El tipo es requerido")
	}
	if err := dto.Tipo.Validate(); err != nil {
		return err
	}

	// Email
	if dto.Email == "" {
		return core.NewValidationError("El email es requerido")
	}
	if validation.Validate(dto.Email, is.Email) != nil {
		return core.NewValidationError("El email es invalido")
	}

	// Telefono
	if dto.Telefono == "" {
		return core.NewValidationError("El telefono es requerido")
	}
	if !regexp.MustCompile(`^0(414|424|412|416|426)[0-9]{7}$`).MatchString(dto.Telefono) {
		return core.NewValidationError("El teléfono debe tener formato 04141234567")
	}

	telefono := dto.Telefono[1:] // "4141234567"

	dto.Telefono = "+58" + telefono

	if validation.Validate(dto.Telefono, is.E164) != nil {
		return core.NewValidationError("Error al formatear el teléfono")
	}

	// Direccion
	if dto.Direccion != nil && *dto.Direccion == "" {
		return core.NewValidationError("La direccion no puede estar vacia")
	}

	return nil
}
