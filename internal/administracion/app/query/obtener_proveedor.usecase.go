package query

import (
	"errors"

	"github.com/Sanaruca/condominio/internal/administracion"
	"github.com/Sanaruca/condominio/internal/core"
	"github.com/Sanaruca/condominio/internal/core/context"
	"github.com/Sanaruca/condominio/internal/core/usecase"
	"gorm.io/gorm"
)

type ObtenerProveedorDTO struct {
	ProveedorID string
}

type ObtenerProveedor usecase.Handler[context.BaseContext, ObtenerProveedorDTO, *administracion.Proveedor]

func NewObtenerProveedor() ObtenerProveedor {
	return &obtenerProveedor{}
}

type obtenerProveedor struct{}

func (uc *obtenerProveedor) Exec(ctx context.BaseContext, dto ObtenerProveedorDTO) (
	*administracion.Proveedor,
	core.Error,
) {

	if err := dto.Validate(); err != nil {
		return nil, err
	}

	proveedor, err := gorm.G[administracion.Proveedor](
		ctx.DB,
	).Where("id = ?", dto.ProveedorID).
		Take(ctx.Context())

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, administracion.ErrProveedorNoEncontrado
	}

	if err != nil {
		return nil, core.WrapError(err)
	}

	return &proveedor, nil

}

func (dto ObtenerProveedorDTO) Validate() core.Error {
	if dto.ProveedorID == "" {
		return core.NewInvalidArgumentError("ProveedorID es requerido")
	}
	return nil
}
