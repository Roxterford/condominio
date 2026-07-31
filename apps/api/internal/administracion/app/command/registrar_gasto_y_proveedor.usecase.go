package command

// Deprecated: Este archivo pertenece al sistema legacy.
// Usar internal/transacciones/app/command/RegistrarTransaccion en su lugar.
import (
	"github.com/Sanaruca/condominio/internal/administracion/models/gasto"
	"github.com/Sanaruca/condominio/internal/core"
	"github.com/Sanaruca/condominio/internal/core/context"
	"github.com/Sanaruca/condominio/internal/core/usecase"
)

type RegistrarGastoYProveedorDTO struct {
	GastoBase
	Proveedor RegistrarProveedorDTO
}

type RegistrarGastoYProveedor usecase.Handler[context.AdminContext, RegistrarGastoYProveedorDTO, *gasto.Gasto]

type registrarGastoYProveedor struct {
	registrarProveedor RegistrarProveedor
	registrarGasto     RegistrarGasto
}

func NewRegistrarGastoYProveedor(
	registrarProveedor RegistrarProveedor,
	registrarGasto RegistrarGasto,
) RegistrarGastoYProveedor {
	if registrarProveedor == nil {
		panic("registrarProveedor is nil")
	}
	if registrarGasto == nil {
		panic("registrarGasto is nil")
	}
	return registrarGastoYProveedor{
		registrarProveedor: registrarProveedor,
		registrarGasto:     registrarGasto,
	}
}

func (uc registrarGastoYProveedor) Exec(
	ctx context.AdminContext,
	input RegistrarGastoYProveedorDTO,
) (*gasto.Gasto, core.Error) {

	proveedor, err := uc.registrarProveedor.Exec(ctx, input.Proveedor)

	if err != nil {
		return nil, err
	}

	_, err = uc.registrarGasto.Exec(ctx, RegistrarGastoDTO{
		GastoBase: input.GastoBase,
		Proveedor: proveedor.ID(),
	})

	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (dto *RegistrarGastoYProveedorDTO) Validate() core.Error {

	err := dto.GastoBase.Validate()

	if err != nil {
		return err
	}

	return dto.Proveedor.Validate()
}
