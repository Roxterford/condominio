package command

import (
	"time"

	"github.com/Sanaruca/condominio/internal/core"
	"github.com/Sanaruca/condominio/internal/core/adapters/ozzo"
	"github.com/Sanaruca/condominio/internal/core/common/events"
	"github.com/Sanaruca/condominio/internal/core/context"
	"github.com/Sanaruca/condominio/internal/core/errors"
	"github.com/Sanaruca/condominio/internal/core/usecase"
	"github.com/Sanaruca/condominio/internal/pagos"
	"github.com/Sanaruca/condominio/internal/pagos/types/metododepago"
	"github.com/Sanaruca/condominio/internal/pagos/types/moneda"
	"github.com/Sanaruca/condominio/internal/services/tasa"
	"github.com/Sanaruca/condominio/internal/villas"
	validation "github.com/go-ozzo/ozzo-validation/v4"
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
	pagos          pagos.PagoRepository
	villas         villas.VillaRepository
	bus_de_eventos events.EventBus
	tasa_service   tasa.TasaService
}

func NewRegistrarPago(
	pago_repository pagos.PagoRepository,
	villa_repository villas.VillaRepository,
	bus_de_eventos events.EventBus,
	tasa_service tasa.TasaService,
) RegistrarPago {
	if pago_repository == nil {
		panic("pago_repository is nil")
	}
	if villa_repository == nil {
		panic("villa_repository is nil")
	}
	if bus_de_eventos == nil {
		panic("bus_de_eventos is nil")
	}
	return &registrarPago{
		pagos:          pago_repository,
		villas:         villa_repository,
		bus_de_eventos: bus_de_eventos,
		tasa_service:   tasa_service,
	}
}

func (uc *registrarPago) Exec(ctx context.AdminContext, input RegistrarPagoDTO) (any, core.Error) {

	if err := input.Validate(); err != nil {
		return nil, err
	}

	if exists, err := uc.villas.Exists(ctx, input.Villa); err != nil {
		return nil, err
	} else if !exists {
		return nil, villas.ErrVillaNoEncontrada
	}

	fechaPago := *input.Fecha
	tasaAUsar := input.Tasa

	if tasaAUsar < 1 && input.Moneda == moneda.VED {
		tasaObtenida, err := uc.tasa_service.ObtenerTasaParaPago(fechaPago)
		if err != nil {
			return nil, errors.New(
				errors.INVALID_ARGUMENT,
				"No se pudo obtener la tasa de tasa. Por favor ingrese la tasa manualmente",
			)
		}
		tasaAUsar = tasaObtenida.Valor
	}

	pago, err := pagos.NuevoPago(
		input.Villa,
		fechaPago,
		input.Metodo,
		input.Monto,
		input.Moneda,
		tasaAUsar,
		input.Referencia,
		ctx.Session().Usuario().ID,
	)
	if err != nil {
		return nil, err
	}

	if err := uc.pagos.Guardar(ctx, pago); err != nil {
		return nil, err
	}

	for _, ev := range pago.PullEvents() {
		// Publicamos en el bus de eventos. Si falla, el Cron Job lo arreglará luego
		if err := uc.bus_de_eventos.Publish(ctx, ev); err != nil {
			// TODO: Logueamos pero no frenamos el proceso, el pago ya es real en la DB
		}
	}

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
		validation.Field(&dto.Tasa, validation.By(func(value any) error {
			if dto.Moneda == moneda.USD {
				return nil
			}
			if value == nil || value.(int) < 1 {
				return errors.New(errors.INVALID_ARGUMENT, "La tasa es requerida para pagos en VED")
			}
			return nil
		})),
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
