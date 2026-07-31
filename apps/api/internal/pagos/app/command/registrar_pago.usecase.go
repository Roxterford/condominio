// Deprecated: Caso de uso legacy. Usar internal/transacciones/app/command/RegistrarTransaccion en su lugar.
package command

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"

	"github.com/Sanaruca/condominio/internal/core"
	"github.com/Sanaruca/condominio/internal/core/adapters/ozzo"
	"github.com/Sanaruca/condominio/internal/core/common/events"
	"github.com/Sanaruca/condominio/internal/core/common/quantity"
	"github.com/Sanaruca/condominio/internal/core/context"
	"github.com/Sanaruca/condominio/internal/core/errors"
	"github.com/Sanaruca/condominio/internal/core/usecase"
	"github.com/Sanaruca/condominio/internal/pagos/models/pago"
	"github.com/Sanaruca/condominio/internal/pagos/types/metododepago"
	"github.com/Sanaruca/condominio/internal/pagos/types/moneda"
	"github.com/Sanaruca/condominio/internal/services/tasa"
	"github.com/Sanaruca/condominio/internal/unidades/models/unidad"
)

type RegistrarPagoDTO struct {
	Unidad     unidad.UnidadCodigo
	Fecha      *time.Time
	Metodo     metododepago.MetodoDePago
	Referencia *string
	Monto      int
	Tasa       int
	Moneda     moneda.Moneda
}

type RegistrarPago usecase.WithContextInput[context.AdminContext, RegistrarPagoDTO]

type registrarPago struct {
	pagoRepo       pago.PagoRepository
	pago_factory   *pago.PagoFactory
	unidades       unidad.UnidadRepository
	bus_de_eventos events.EventBus
	tasa_service   tasa.TasaService

	qf *quantity.QuantityFactory
}

func NewRegistrarPago(
	pago_repository pago.PagoRepository,
	pago_factory *pago.PagoFactory,
	unidad_repository unidad.UnidadRepository,
	bus_de_eventos events.EventBus,
	tasa_service tasa.TasaService,
	quantity_factory *quantity.QuantityFactory,
) RegistrarPago {
	if pago_repository == nil {
		panic("pago_repository is nil")
	}
	if pago_factory == nil {
		panic("pago_factory is nil")
	}
	if unidad_repository == nil {
		panic("unidad_repository is nil")
	}
	if bus_de_eventos == nil {
		panic("bus_de_eventos is nil")
	}
	if quantity_factory == nil {
		panic("qf is nil")
	}
	return &registrarPago{
		pagoRepo:       pago_repository,
		unidades:       unidad_repository,
		bus_de_eventos: bus_de_eventos,
		tasa_service:   tasa_service,
		qf:             quantity_factory,
		pago_factory:   pago_factory,
	}
}

func (uc *registrarPago) Exec(ctx context.AdminContext, input RegistrarPagoDTO) (any, core.Error) {

	if err := input.Validate(); err != nil {
		return nil, err
	}

	if exists, err := uc.unidades.ExistsCodigo(ctx, input.Unidad); err != nil {
		return nil, err
	} else if !exists {
		return nil, unidad.ErrUnidadNoEncontrada
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
		tasaAUsar = int(tasaObtenida.Valor.Value()) // TODO: check
	}

	_pago, err := uc.pago_factory.Nuevo(
		input.Unidad,
		fechaPago,
		input.Metodo,
		uc.qf.Assemble(int64(input.Monto)),
		input.Moneda,
		uc.qf.Assemble(int64(tasaAUsar)),
		input.Referencia,
		ctx.Session().Usuario().ID,
	)
	if err != nil {
		return nil, err
	}

	if err := uc.pagoRepo.Guardar(ctx, _pago); err != nil {
		return nil, err
	}

	for _, ev := range _pago.PullEvents() {
		// Publicamos en el bus de eventos. Si falla, el Cron Job lo arreglará luego
		if err := uc.bus_de_eventos.Publish(ctx, ev); err != nil {
			// TODO: Logueamos pero no frenamos el proceso, el pago ya es real en la DB
		}
	}

	return _pago, nil
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
		validation.Field(&dto.Unidad, validation.Required),
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
