package command

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"

	"github.com/Sanaruca/condominio/internal/administracion/models/cuota"
	"github.com/Sanaruca/condominio/internal/core"
	"github.com/Sanaruca/condominio/internal/core/adapters/ozzo"
	"github.com/Sanaruca/condominio/internal/core/common"
	"github.com/Sanaruca/condominio/internal/core/common/filter"
	"github.com/Sanaruca/condominio/internal/core/common/mes"
	cc "github.com/Sanaruca/condominio/internal/core/context"
	"github.com/Sanaruca/condominio/internal/core/usecase"
	"github.com/Sanaruca/condominio/internal/finanzas/models/operacion"
)

type RegistrarCuotaRegularDTO struct {
	Mes         mes.Mes
	Anio        int
	FechaLimite time.Time
	Gastos      []string
}

type CuotaUoWDeps struct {
	Cuotas      cuota.CuotaRepository
	Operaciones operacion.OperacionRepository
}

type RegistrarCuotaRegular usecase.Handler[cc.AdminContext, RegistrarCuotaRegularDTO, *cuota.CuotaRegular]

type registrarCuotaRegular struct {
	cuotas      cuota.CuotaRepository
	operaciones operacion.OperacionRepository
	cf          *cuota.CuotaFactory
	uow         common.UnitOfWork[CuotaUoWDeps]
}

func NewRegistrarCuotaRegular(
	cuotaRepository cuota.CuotaRepository,
	operacionRepository operacion.OperacionRepository,
	cuotaFactory *cuota.CuotaFactory,
	uow common.UnitOfWork[CuotaUoWDeps],
) RegistrarCuotaRegular {

	if cuotaRepository == nil {
		panic("cuotaRepository is nil")
	}

	if operacionRepository == nil {
		panic("operacionRepository is nil")
	}

	if cuotaFactory == nil {
		panic("cuotaFactory is nil")
	}

	if uow == nil {
		panic("uow is nil")
	}

	return &registrarCuotaRegular{
		cuotas:      cuotaRepository,
		operaciones: operacionRepository,
		cf:          cuotaFactory,
		uow:         uow,
	}
}

func (uc registrarCuotaRegular) Exec(
	ctx cc.AdminContext,
	input RegistrarCuotaRegularDTO,
) (*cuota.CuotaRegular, core.Error) {

	if len(input.Gastos) == 0 {
		return nil, core.NewValidationError("Debe proporcionar al menos un gasto")
	}

	ftr, err := filter.Parse(map[string]any{
		"id": map[string]any{
			"in": input.Gastos,
		},
	})
	if err != nil {
		return nil, core.WrapError(err)
	}

	gastos, err := uc.operaciones.Obtener(ctx, ftr, common.Paginator{})
	if err != nil {
		return nil, core.WrapError(err)
	}

	if len(gastos.Data) != len(input.Gastos) {
		return nil, core.NewValidationError(
			"Uno o varios de los gastos proporcionados no fue encontrado",
		)
	}

	monto := gastos.Data[0].Total()
	for _, gasto := range gastos.Data[1:] {
		monto = monto.HappyAdd(gasto.Total())
	}

	_cuota, err := uc.cf.NuevaRegular(
		int(monto.Value()),
		input.Mes,
		input.Anio,
		ctx.Session().Usuario().ID,
	)
	if err != nil {
		return nil, core.WrapError(err)
	}

	txerr := uc.uow.Do(ctx, func(deps CuotaUoWDeps) error {

		if _, err := deps.Cuotas.Guardar(ctx, _cuota); err != nil {
			return err
		}

		for _, gasto := range gastos.Data {
			if err := gasto.AsignarCuota(_cuota.ID().String()); err != nil {
				return err
			}
			if err := deps.Operaciones.Guardar(ctx, &gasto); err != nil {
				return err
			}
		}

		return nil
	})

	if txerr != nil {
		return nil, core.WrapError(txerr)
	}

	return _cuota, nil
}

func (input *RegistrarCuotaRegularDTO) Validate() core.Error {

	unduplicatedIDs := core.NewSetFromSlice(input.Gastos, func(it string) string { return it })
	input.Gastos = unduplicatedIDs.ToSlice()

	if err := input.Mes.Validate(); err != nil {
		return err
	}

	err := validation.ValidateStruct(input,
		validation.Field(
			&input.Gastos,
			validation.Required,
			validation.Max(100),
		),
	)

	if err != nil {
		return ozzo.FirstOzzoErrorAdapter(input, err)
	}

	return nil
}
