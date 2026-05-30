package command

import (
	"slices"
	"strconv"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"

	"github.com/Sanaruca/condominio/internal/administracion/models/cuota"
	"github.com/Sanaruca/condominio/internal/administracion/models/gasto"
	"github.com/Sanaruca/condominio/internal/core"
	"github.com/Sanaruca/condominio/internal/core/adapters/ozzo"
	"github.com/Sanaruca/condominio/internal/core/common"
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
	Gastos []gasto.GastoID `json:"gastos"`

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

type RegistrarCuotaDeps struct {
	Cuotas cuota.CuotaRepository
	Gastos gasto.GastoRepository
}

type registrarCuota struct {
	cuotaFactory *cuota.CuotaFactory

	uow common.UnitOfWork[RegistrarCuotaDeps]
}

func NewRegistrarCuota(
	cuotaFactory *cuota.CuotaFactory,
	uow common.UnitOfWork[RegistrarCuotaDeps],
) RegistrarCuota {

	if cuotaFactory == nil {
		panic("cuotaFactory is nil")
	}

	return &registrarCuota{
		cuotaFactory: cuotaFactory,
		uow:          uow,
	}
}

func (uc *registrarCuota) Exec(
	ctx cc.AdminContext,
	input RegistrarCuotaDTO,
) (cuota.Cuota, core.Error) {

	if err := input.Validate(); err != nil {
		return nil, err
	}

	var _cuota cuota.Cuota

	// Ejecutar el bloque transaccional usando la interfaz de UnitOfWork
	err := uc.uow.Do(ctx, func(tx RegistrarCuotaDeps) error {

		gastos, err := tx.Gastos.ObtenerPorIDs(ctx, input.Gastos)

		if err != nil {
			return err
		}

		var monto_total int

		if len(gastos) != len(input.Gastos) {
			encontradosIDs := make(map[gasto.GastoID]struct{}, len(gastos))
			for _, g := range gastos {
				encontradosIDs[g.ID()] = struct{}{}
				monto_total += int(g.Monto().Value())
			}

			for _, id := range input.Gastos {
				if _, ok := encontradosIDs[id]; !ok {
					// TODO: deberia ser un notfound error
					return core.NewValidationError(
						"Gasto no encontrado: '" + string(id) + "'",
					)
				}
			}
		} else {
			for _, g := range gastos {
				monto_total += int(g.Monto().Value())
			}
		}
		registrador := ctx.Session().Usuario().ID

		if input.Tipo == TipoCuotaRegular {
			_cuota, err = uc.cuotaFactory.NuevaRegular(
				monto_total,
				input.Mes,
				input.Anio,
				registrador,
			)
		} else {
			_cuota, err = uc.cuotaFactory.NuevaEspecial(
				input.Mes,
				input.Anio,
				monto_total,
				input.Titulo,
				input.Descripcion,
				input.Justificacion,
				*input.FechaLimite,
				0,
				registrador,
			)
		}

		if err != nil {
			return err
		}

		_, err = tx.Cuotas.Guardar(ctx, _cuota)
		if err != nil {
			return err
		}

		cuotaID := _cuota.ID().String()

		for i := range gastos {
			// TODO: !Esto es tan importante que quisa se debamos mesclar la
			// enitdad cuota con gastos a pesar de que cargemos muchos gastos en
			// el modelo cuota sin usar
			if err := gastos[i].SetCuota(cuotaID); err != nil {
				return err
			}
			if err := tx.Gastos.Actualizar(ctx, gastos[i]); err != nil {
				return err
			}
		}

		return nil

	})

	return _cuota, core.WrapError(err)
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

	dto.Gastos = slices.Compact(dto.Gastos)

	err := validation.ValidateStruct(
		dto,
		validation.Field(
			&dto.Gastos,
			validation.Required,
			validation.Length(1, 100),
			validation.Each(validation.Required),
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
