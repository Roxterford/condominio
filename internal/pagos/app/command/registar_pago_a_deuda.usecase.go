package command

import (
	"errors"
	"fmt"
	"time"

	"github.com/Sanaruca/condominio/internal/core"
	"github.com/Sanaruca/condominio/internal/core/adapters/ozzo"
	coreError "github.com/Sanaruca/condominio/internal/core/errors"

	"github.com/Sanaruca/condominio/internal/core/context"
	"github.com/Sanaruca/condominio/internal/core/usecase"
	"github.com/Sanaruca/condominio/internal/pagos"
	"github.com/Sanaruca/condominio/internal/villas"
	"github.com/Sanaruca/condominio/internal/villas/types/estadodeuda"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/lucsky/cuid"
	"gorm.io/gorm"
)

type RegistarPagoADeudaDTO struct {
	PagoID string
}

type RegistarPagoADeuda usecase.WithContextInput[context.BaseContext, RegistarPagoADeudaDTO]

type registarPagoADeuda struct{}

func NewRegistarPagoADeuda() RegistarPagoADeuda {
	return &registarPagoADeuda{}
}

func (uc *registarPagoADeuda) Exec(
	ctx context.BaseContext,
	dto RegistarPagoADeudaDTO,
) (any, core.Error) {

	if err := dto.Validate(); err != nil {
		return nil, err
	}

	pago, err := gorm.G[pagos.Pago](
		ctx.DB,
	).Where("id = ?", dto.PagoID).
		First(ctx)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, pagos.ErrPagoNoEncontrado
	}

	if err != nil {
		return nil, coreError.Wrap(err)
	}

	deuda, err := gorm.G[villas.Deuda](ctx.DB).Where(
		"villa = ? AND estado <> ?",
		pago.Villa,
		estadodeuda.Pagada,
	).Select("id", "deuda").Order("registro asc").First(ctx)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if err != nil {
		return nil, coreError.Wrap(err)
	}

	fmt.Println(pago.Total)
	fmt.Println(pago.Cuenta)
	if pago.Cuenta <= 0 {
		return nil, nil
	}

	var destinado int

	if pago.Cuenta >= deuda.Deuda {
		destinado = deuda.Deuda
	} else {
		destinado = pago.Cuenta
	}

	fmt.Println(destinado)

	destino := pagos.DestinoDePago{
		ID:        cuid.New(),
		Pago:      pago.ID,
		Deuda:     deuda.ID,
		Destinado: destinado,
		Fecha:     time.Now(),
	}

	if err := gorm.G[pagos.DestinoDePago](ctx.DB).Create(ctx, &destino); err != nil {
		return nil, coreError.Wrap(err)
	}

	return nil, nil
}

func (dto RegistarPagoADeudaDTO) Validate() core.Error {

	err := validation.ValidateStruct(&dto,
		validation.Field(&dto.PagoID,
			validation.Required,
		),
	)

	return ozzo.FirstOzzoErrorAdapter(dto, err)
}
