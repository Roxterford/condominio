package command

import (
	"strconv"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"

	"github.com/Sanaruca/condominio/internal/administracion/models/cuota"
	"github.com/Sanaruca/condominio/internal/core"
	"github.com/Sanaruca/condominio/internal/core/adapters/ozzo"
	"github.com/Sanaruca/condominio/internal/core/common"
	"github.com/Sanaruca/condominio/internal/core/common/filter"
	"github.com/Sanaruca/condominio/internal/core/common/mes"
	cc "github.com/Sanaruca/condominio/internal/core/context"
	"github.com/Sanaruca/condominio/internal/core/testing/utils"
	"github.com/Sanaruca/condominio/internal/core/usecase"
	"github.com/Sanaruca/condominio/internal/finanzas/models/transaccion"
	"github.com/Sanaruca/condominio/internal/finanzas/types/tipodemovimiento"
)

type TipoDeCuota string

const (
	TipoCuotaRegular  TipoDeCuota = "regular"
	TipoCuotaEspecial TipoDeCuota = "especial"
)

type RegistrarCuotaDTO struct {
	Gastos []string

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
	cuotaFactory  *cuota.CuotaFactory
	cuotas        cuota.CuotaRepository
	transacciones transaccion.TransaccionRepository
}

func NewRegistrarCuota(
	cuotaFactory *cuota.CuotaFactory,
	cuotaRepository cuota.CuotaRepository,
	transaccionRepository transaccion.TransaccionRepository,
) RegistrarCuota {

	if cuotaFactory == nil {
		panic("cuotaFactory is nil")
	}

	if cuotaRepository == nil {
		panic("cuotaRepository is nil")
	}

	if transaccionRepository == nil {
		panic("transaccionRepository is nil")
	}

	return &registrarCuota{
		cuotaFactory:  cuotaFactory,
		cuotas:        cuotaRepository,
		transacciones: transaccionRepository,
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
	ftr, _ := filter.Parse(map[string]any{
		"id": map[string]any{
			"in": input.Gastos,
		},
	})

	gastos, err := uc.transacciones.Obtener(
		ctx,
		ftr,
		common.Paginator{Limit: 100},
		utils.Ptr(tipodemovimiento.Debito),
	)

	if err != nil {
		return nil, err
	}

	var monto int

	for _, t := range gastos.Data {
		monto += int(t.TotalDebitos().Value())
	}

	var _cuota cuota.Cuota

	if input.Tipo == TipoCuotaRegular {
		_cuota, err = uc.cuotaFactory.NuevaRegular(
			monto,
			input.Mes,
			input.Anio,
			registrador,
		)
	} else {
		_cuota, err = uc.cuotaFactory.NuevaEspecial(
			input.Mes,
			input.Anio,
			monto,
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

func (input *RegistrarCuotaDTO) Validate() core.Error {

	now := time.Now()
	if input.Mes == 0 {
		input.Mes = mes.Mes(now.Month() + 1)
	}
	if input.Anio == 0 {
		input.Anio = now.Year()
	}
	if input.FechaLimite == nil {
		fechaLimite := now.AddDate(0, 1, 0) // TODO: hacer configurable
		input.FechaLimite = &fechaLimite
	}

	err := validation.ValidateStruct(
		input,
		validation.Field(
			&input.Gastos,
			validation.Required,
			validation.Length(1, 0),
		),
		validation.Field(
			&input.Tipo,
			validation.Required,
		),
	)

	if err != nil {
		return ozzo.FirstOzzoErrorAdapter(input, err)
	}

	// TODO: ver
	if len(input.Gastos) > 100 {
		return core.NewValidationError("no se pueden procesar mas de 100 gastos por el momento")
	}

	if input.Tipo == TipoCuotaRegular {
		if err := input.Mes.Validate(); err != nil {
			return err
		}

		if input.Anio < 2000 {
			return core.NewValidationError("el año debe ser mayor a 2000")
		}

		if input.Titulo == "" {
			input.Titulo = "Mensualidad " + input.Mes.String() + ", " + strconv.Itoa(input.Anio)
		}
	}

	if input.Tipo == TipoCuotaEspecial {
		if input.Titulo == "" {
			input.Titulo = "Cuota Especial - " + input.Mes.String() + ", " + strconv.Itoa(
				input.Anio,
			)
		}

		if err := validation.ValidateStruct(
			input,
			validation.Field(&input.Descripcion, validation.Required, validation.Length(1, 500)),
			validation.Field(&input.Justificacion, validation.Required, validation.Length(1, 500)),
		); err != nil {
			return ozzo.FirstOzzoErrorAdapter(input, err)
		}
	}

	return nil
}
