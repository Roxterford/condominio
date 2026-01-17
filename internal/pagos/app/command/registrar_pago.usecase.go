package command

import (
	"time"

	"github.com/Sanaruca/condominio/internal/core"
	"github.com/Sanaruca/condominio/internal/core/adapters/ozzo"
	"github.com/Sanaruca/condominio/internal/core/context"
	"github.com/Sanaruca/condominio/internal/core/errors"
	"github.com/Sanaruca/condominio/internal/core/usecase"
	"github.com/Sanaruca/condominio/internal/pagos"
	"github.com/Sanaruca/condominio/internal/pagos/types/metododepago"
	"github.com/Sanaruca/condominio/internal/pagos/types/moneda"
	"github.com/Sanaruca/condominio/internal/villas"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/lucsky/cuid"
	"gorm.io/gorm"
)

type RegistrarPagoDTO struct {
	Villa      int
	Fecha      *time.Time
	Metodo     metododepago.MetodoDePago
	Referencia *string
	Monto      int
	Tasa       int
	Moneda     moneda.Moneda
}

type RegistrarPago usecase.WithContextInput[context.AdminContext, RegistrarPagoDTO]

type registrarPago struct {
	RegistrarPagoADeuda RegistarPagoADeuda
}

func NewRegistrarPago() RegistrarPago {
	return &registrarPago{}
}

func (uc *registrarPago) Exec(ctx context.AdminContext, dto RegistrarPagoDTO) (any, core.Error) {

	if err := dto.Validate(); err != nil {
		return nil, err
	}

	_, err := gorm.G[villas.Villa](ctx.DB).Where("numero = ?", dto.Villa).Select("id").First(ctx)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, villas.ErrVillaNoEncontrada
	}

	pago := pagos.IPago{
		ID:             cuid.New(),
		Villa:          dto.Villa,
		Fecha:          time.Time{},
		Metodo:         dto.Metodo,
		Monto:          dto.Monto,
		Referencia:     dto.Referencia,
		Moneda:         dto.Moneda,
		RegistradoPor:  ctx.Session.Usuario.ID,
		Registro:       time.Time{},
		ActualizadoPor: ctx.Session.Usuario.ID,
		Actualizacion:  time.Time{},
		Tasa:           dto.Tasa,
	}

	gorm.G[pagos.IPago](ctx.DB).Create(ctx, &pago)

	go uc.RegistrarPagoADeuda.Exec(ctx.BaseContext, RegistarPagoADeudaDTO{pago.ID})

	return pago, nil
}

func (dto *RegistrarPagoDTO) Validate() core.Error {

	if err := dto.Moneda.Validate(); err != nil {
		return err
	}
	if err := dto.Metodo.Validate(); err != nil {
		return err
	}

	err := validation.ValidateStruct(
		dto,
		validation.Field(&dto.Villa, validation.Required, validation.Min(1)),
		validation.Field(&dto.Tasa, validation.Required, validation.Min(1)),
		validation.Field(&dto.Monto, validation.Required, validation.Min(1)),
		validation.Field(&dto.Referencia, validation.NilOrNotEmpty, validation.Length(1, 50)),
	)

	if dto.Fecha != nil {
		// No permitir fechas futuras
		if dto.Fecha.After(time.Now()) {
			return core.NewError(errors.INVALID_ARGUMENT, "La fecha no puede ser futura")
		}
	} else {
		now := time.Now()
		dto.Fecha = &now
	}

	return ozzo.FirstOzzoErrorAdapter(dto, err)
}
