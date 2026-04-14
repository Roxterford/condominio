package command

import (
	"github.com/Sanaruca/condominio/internal/administracion/models/cuota"
	"github.com/Sanaruca/condominio/internal/administracion/models/deuda"
	"github.com/Sanaruca/condominio/internal/core"
	"github.com/Sanaruca/condominio/internal/core/adapters/ozzo"
	cc "github.com/Sanaruca/condominio/internal/core/context"
	"github.com/Sanaruca/condominio/internal/core/usecase"
	"github.com/Sanaruca/condominio/internal/villas/models/villa"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type AplicarCuotaDTO struct {
	CuotaID string
}

type AplicarCuota usecase.Handler[cc.BaseContext, AplicarCuotaDTO, any]

type aplicarCuota struct {
	cuotaRepo         cuota.CuotaRepository
	villaRepo         villa.VillaRepository
	deudaRepo         deuda.DeudaRepository
	facturacionPolicy villa.FacturacionPolicy
	deudaFactory      *deuda.DeudaFactory
}

func NewAplicaCuota(
	cuotaRepo cuota.CuotaRepository,
	villaRepo villa.VillaRepository,
	deudaRepo deuda.DeudaRepository,
	policy villa.FacturacionPolicy,
	deudaFactory *deuda.DeudaFactory,
) AplicarCuota {
	if cuotaRepo == nil {
		panic("cuotaRepo is nil")
	}
	if villaRepo == nil {
		panic("villaRepo is nil")
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
		villaRepo:         villaRepo,
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

	villas, err := uc.villaRepo.ObtenerTodas(ctx)
	if err != nil {
		return nil, err
	}

	montoCuota := int(_cuota.Monto().Value())
	villasConDeuda := 0

	for _, villaNumero := range villas {
		estado, err := uc.villaRepo.ObtenerEstado(ctx, villaNumero)
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
			villaNumero,
			montoCuota,
		)
		if err != nil {
			return nil, err
		}

		if err := uc.deudaRepo.Guardar(ctx, nuevaDeuda); err != nil {
			return nil, err
		}
		villasConDeuda++
	}

	return map[string]int{
		"villas_procesadas": villasConDeuda,
	}, nil
}

func (dto AplicarCuotaDTO) Validate() core.Error {
	err := validation.ValidateStruct(&dto,
		validation.Field(&dto.CuotaID, validation.Required),
	)
	return ozzo.FirstOzzoErrorAdapter(dto, err)
}
