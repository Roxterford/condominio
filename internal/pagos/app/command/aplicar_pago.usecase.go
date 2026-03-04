package command

import (
	"github.com/Sanaruca/condominio/internal/core"
	"github.com/Sanaruca/condominio/internal/core/adapters/ozzo"

	"github.com/Sanaruca/condominio/internal/core/context"
	"github.com/Sanaruca/condominio/internal/core/usecase"
	"github.com/Sanaruca/condominio/internal/pagos"
	"github.com/Sanaruca/condominio/internal/villas"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type AplicarPagoDTO struct {
	PagoID string
}

// RegistarPagoADeuda se encarga de registrar un pago a las deudas pendientes de una villa segun la deuda mas antigua.
type AplicarPago usecase.WithContextInput[context.BaseContext, AplicarPagoDTO]

type aplicarPago struct {
	pagos  pagos.PagoRepository
	villas villas.VillaRepository
}

func NewAplicarPago() AplicarPago {
	return &aplicarPago{}
}

func (uc *aplicarPago) Exec(
	ctx context.BaseContext,
	dto AplicarPagoDTO,
) (any, core.Error) {

	if err := dto.Validate(); err != nil {
		return nil, err
	}

	pago, err := uc.pagos.GetByID(ctx, dto.PagoID)

	if err != nil {
		return nil, err
	}

	deuda, err := uc.villas.GetLastDeudaWhereNotPagada(ctx, pago.Villa())

	if err != nil {
		return nil, err
	}

	// Si no hay deudas pendientes, no hay nada que aplicar
	if deuda == nil {
		return nil, nil
	}

	_, err = pago.AplicarPagoADeuda(deuda.ID(), deuda.Restante())
	if err != nil {
		return nil, err
	}

	if err := uc.pagos.Guardar(ctx, pago); err != nil {
		return nil, err
	}

	return nil, nil
}

func (dto AplicarPagoDTO) Validate() core.Error {

	err := validation.ValidateStruct(&dto,
		validation.Field(&dto.PagoID,
			validation.Required,
		),
	)

	return ozzo.FirstOzzoErrorAdapter(dto, err)
}
