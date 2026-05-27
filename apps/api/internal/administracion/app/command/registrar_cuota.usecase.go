package command

import (
	"strconv"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"

	"github.com/Sanaruca/condominio/internal/administracion/models/cuota"
	"github.com/Sanaruca/condominio/internal/administracion/models/gasto"
	"github.com/Sanaruca/condominio/internal/core"
	"github.com/Sanaruca/condominio/internal/core/adapters/ozzo"
	"github.com/Sanaruca/condominio/internal/core/common/events"
	"github.com/Sanaruca/condominio/internal/core/common/mes"
	"github.com/Sanaruca/condominio/internal/core/context"
	"github.com/Sanaruca/condominio/internal/core/usecase"
)

type TipoDeCuota string

const (
	TipoCuotaRegular  TipoDeCuota = "regular"
	TipoCuotaEspecial TipoDeCuota = "especial"
)

type RegistrarCuotaDTO struct {
	Gastos []gasto.GastoID

	Tipo TipoDeCuota

	// Para regular
	Mes  mes.Mes
	Anio int

	FechaLimite *time.Time

	// Para especial
	Titulo        string
	Descripcion   string
	Justificacion string
}

type RegistrarCuota usecase.WithContextInput[context.AdminContext, RegistrarCuotaDTO]

type registrarCuota struct {
	gastoRepo    gasto.GastoRepository
	cuotaRepo    cuota.CuotaRepository
	cuotaFactory *cuota.CuotaFactory
	eventBus     events.EventBus
}

func NewRegistrarCuota(
	gastoRepo gasto.GastoRepository,
	cuotaRepo cuota.CuotaRepository,
	cuotaFactory *cuota.CuotaFactory,
	eventBus events.EventBus,
) RegistrarCuota {
	if gastoRepo == nil {
		panic("gastoRepo is nil")
	}
	if cuotaRepo == nil {
		panic("cuotaRepo is nil")
	}
	if cuotaFactory == nil {
		panic("cuotaFactory is nil")
	}
	if eventBus == nil {
		panic("eventBus is nil")
	}
	return &registrarCuota{
		gastoRepo:    gastoRepo,
		cuotaRepo:    cuotaRepo,
		cuotaFactory: cuotaFactory,
		eventBus:     eventBus,
	}
}

func (uc *registrarCuota) Exec(
	ctx context.AdminContext,
	input RegistrarCuotaDTO,
) (any, core.Error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	gastoIDs := make([]gasto.GastoID, len(input.Gastos))
	for i, id := range input.Gastos {
		gastoIDs[i] = gasto.GastoID(id)
	}

	gastos, err := uc.gastoRepo.ObtenerPorIDs(ctx, gastoIDs)
	if err != nil {
		return nil, err
	}

	if len(gastos) == 0 {
		return nil, core.NewValidationError("debe incluir al menos un gasto")
	}

	var totalUSD int
	for _, g := range gastos {
		totalUSD += int(g.Total().Value()) // TODO: check
	}

	var _cuota cuota.Cuota
	registrador := ctx.Session().Usuario().ID

	if input.Tipo == TipoCuotaRegular {
		_cuota, err = uc.cuotaFactory.NuevaRegular(
			totalUSD,
			input.Mes,
			input.Anio,
			registrador,
		)
	} else {
		_cuota, err = uc.cuotaFactory.NuevaEspecial(
			input.Mes,
			input.Anio,
			totalUSD,
			input.Titulo,
			input.Descripcion,
			input.Justificacion,
			*input.FechaLimite,
			0,
			registrador,
		)
	}

	if err != nil {
		return nil, err
	}

	_, err = uc.cuotaRepo.Guardar(ctx, _cuota)
	if err != nil {
		return nil, err
	}

	cuotaID := _cuota.ID()
	cuotaIDStr := string(cuotaID)

	for i := range gastos {
		gastos[i].SetCuota(cuotaIDStr)
		if err := uc.gastoRepo.Actualizar(ctx, gastos[i]); err != nil {
			return nil, err
		}
	}

	for _, ev := range _cuota.(interface{ PullEvents() []events.Event }).PullEvents() {
		if err := uc.eventBus.Publish(ctx, ev); err != nil {
			// Loggear pero continuar
		}
	}

	return _cuota, nil
}

func (dto *RegistrarCuotaDTO) Validate() core.Error {

	now := time.Now()
	if dto.Mes == 0 {
		dto.Mes = mes.Mes(now.Month() + 1)
	}
	if dto.Anio == 0 {
		dto.Anio = now.Year()
	}
	if dto.FechaLimite == nil {
		fechaLimite := now.AddDate(0, 1, 0) // TODO: hacer configurable
		dto.FechaLimite = &fechaLimite
	}

	err := validation.ValidateStruct(
		dto,
		validation.Field(&dto.Gastos, validation.Required, validation.Length(1, 100)),
		validation.Field(
			&dto.Tipo,
			validation.Required,
			validation.In(TipoCuotaRegular, TipoCuotaEspecial),
		),
	)

	if err != nil {
		return ozzo.FirstOzzoErrorAdapter(dto, err)
	}

	if dto.Tipo == TipoCuotaRegular {
		if err := dto.Mes.Validate(); err != nil {
			return err
		}

		if dto.Anio < 2000 {
			return core.NewValidationError("el año debe ser mayor a 2000")
		}

		if dto.Titulo == "" {
			dto.Titulo = "Mensualidad " + dto.Mes.String() + ", " + strconv.Itoa(dto.Anio)
		}
	}

	if dto.Tipo == TipoCuotaEspecial {
		if dto.Titulo == "" {
			dto.Titulo = "Cuota Especial - " + dto.Mes.String() + ", " + strconv.Itoa(dto.Anio)
		}

		if err := validation.ValidateStruct(
			dto,
			validation.Field(&dto.Descripcion, validation.Required, validation.Length(1, 500)),
			validation.Field(&dto.Justificacion, validation.Required, validation.Length(1, 500)),
		); err != nil {
			return ozzo.FirstOzzoErrorAdapter(dto, err)
		}
	}

	return nil
}
