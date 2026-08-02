package command

import (
	"strconv"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"

	"github.com/Sanaruca/condominio/internal/administracion/models/cuota"
	"github.com/Sanaruca/condominio/internal/core"
	"github.com/Sanaruca/condominio/internal/core/adapters/ozzo"
	"github.com/Sanaruca/condominio/internal/core/common/mes"
	cc "github.com/Sanaruca/condominio/internal/core/context"
	"github.com/Sanaruca/condominio/internal/core/usecase"
)

type TipoDeCuota string

const (
	TipoCuotaRegular  TipoDeCuota = "regular"
	TipoCuotaEspecial TipoDeCuota = "especial"
)

type RegistrarCuotaDTO struct {
	MontoTotal int `json:"monto_total"`

	Tipo TipoDeCuota `json:"tipo"`

	// Para regular
	Mes  mes.Mes `json:"mes"`
	Anio int     `json:"anio"`

	FechaLimite *time.Time `json:"fecha_limite"`
	// Para especial
	Titulo        string `json:"titulo"`
	Descripcion   string `json:"descripcion"`
	Justificacion string `json:"justificacion"`
}

type RegistrarCuota usecase.Handler[cc.AdminContext, RegistrarCuotaDTO, cuota.Cuota]

type registrarCuota struct {
	cuotaFactory *cuota.CuotaFactory
	cuotas       cuota.CuotaRepository
}

func NewRegistrarCuota(
	cuotaFactory *cuota.CuotaFactory,
	cuotaRepository cuota.CuotaRepository,
) RegistrarCuota {

	if cuotaFactory == nil {
		panic("cuotaFactory is nil")
	}

	if cuotaRepository == nil {
		panic("cuotaRepository is nil")
	}

	return &registrarCuota{
		cuotaFactory: cuotaFactory,
		cuotas:       cuotaRepository,
	}
}

func (uc *registrarCuota) Exec(
	ctx cc.AdminContext,
	input RegistrarCuotaDTO,
) (cuota.Cuota, core.Error) {

	if err := input.Validate(); err != nil {
		return nil, err
	}

	registrador := ctx.Session().Usuario().ID

	var _cuota cuota.Cuota
	var err core.Error

	if input.Tipo == TipoCuotaRegular {
		_cuota, err = uc.cuotaFactory.NuevaRegular(
			input.MontoTotal,
			input.Mes,
			input.Anio,
			registrador,
		)
	} else {
		_cuota, err = uc.cuotaFactory.NuevaEspecial(
			input.Mes,
			input.Anio,
			input.MontoTotal,
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

	if _, err := uc.cuotas.Guardar(ctx, _cuota); err != nil {
		return nil, err
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
		validation.Field(
			&dto.MontoTotal,
			validation.Required,
			validation.Min(1),
		),
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
