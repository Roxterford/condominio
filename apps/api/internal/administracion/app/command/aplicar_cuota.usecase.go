package command

import (
	"github.com/Sanaruca/condominio/internal/administracion/models/cuota"
	"github.com/Sanaruca/condominio/internal/administracion/models/deuda"
	"github.com/Sanaruca/condominio/internal/core"
	"github.com/Sanaruca/condominio/internal/core/adapters/ozzo"
	cc "github.com/Sanaruca/condominio/internal/core/context"
	"github.com/Sanaruca/condominio/internal/core/usecase"
	"github.com/Sanaruca/condominio/internal/unidades/models/unidad"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type AplicarCuotaDTO struct {
	CuotaID string
}

type AplicarCuota usecase.Handler[cc.BaseContext, AplicarCuotaDTO, any]

type aplicarCuota struct {
	cuotaRepo         cuota.CuotaRepository
	unidadRepo        unidad.UnidadRepository
	deudaRepo         deuda.DeudaRepository
	facturacionPolicy unidad.FacturacionPolicy
	deudaFactory      *deuda.DeudaFactory
}

func NewAplicarCuota(
	cuotaRepo cuota.CuotaRepository,
	unidadRepo unidad.UnidadRepository,
	deudaRepo deuda.DeudaRepository,
	policy unidad.FacturacionPolicy,
	deudaFactory *deuda.DeudaFactory,
) AplicarCuota {
	if cuotaRepo == nil {
		panic("cuotaRepo is nil")
	}
	if unidadRepo == nil {
		panic("unidadRepo is nil")
	}
	if deudaRepo == nil {
		panic("deudaRepo is nil")
	}
	if policy == nil {
		panic("policy is missing")
	}
	if deudaFactory == nil {
		panic("deudaFactory is nil")
	}
	return aplicarCuota{
		cuotaRepo:         cuotaRepo,
		unidadRepo:        unidadRepo,
		deudaRepo:         deudaRepo,
		facturacionPolicy: policy,
		deudaFactory:      deudaFactory,
	}
}

func (uc aplicarCuota) Exec(ctx cc.BaseContext, input AplicarCuotaDTO) (any, core.Error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	_cuota, err := uc.cuotaRepo.ObtenerPorID(ctx, cuota.CuotaID(input.CuotaID))
	if err != nil {
		return nil, err
	}
	if _cuota == nil {
		return nil, cuota.ErrCuotaNoEncontrada
	}

	unidades, err := uc.unidadRepo.ObtenerTodas(ctx)
	if err != nil {
		return nil, err
	}

	montoCuota := int(_cuota.Monto().Value())
	unidadesConDeuda := 0

	for _, _unidad := range unidades {
		estado, err := uc.unidadRepo.ObtenerEstado(ctx, _unidad.Codigo())
		if err != nil {
			return nil, err
		}

		var generaDeuda bool
		if _cuota.AsRegular() != nil {
			generaDeuda = uc.facturacionPolicy.GeneraDeudaPorCuotaRegular(estado)
		} else if _cuota.AsEspecial() != nil {
			generaDeuda = uc.facturacionPolicy.GeneraDeudaPorCuotaEspecial(estado)
		}

		if !generaDeuda {
			continue
		}

		nuevaDeuda, err := uc.deudaFactory.NuevaDeuda(
			_cuota.ID(),
			_unidad,
			montoCuota,
		)
		if err != nil {
			return nil, err
		}

		if err := uc.deudaRepo.Guardar(ctx, nuevaDeuda); err != nil {
			return nil, err
		}
		unidadesConDeuda++
	}

	return map[string]int{
		"unidades_procesadas": unidadesConDeuda,
	}, nil
}

func (dto AplicarCuotaDTO) Validate() core.Error {
	err := validation.ValidateStruct(&dto,
		validation.Field(&dto.CuotaID, validation.Required),
	)
	return ozzo.FirstOzzoErrorAdapter(dto, err)
}
