package command

import (
	"time"

	"github.com/Sanaruca/condominio/internal/administracion/models/gasto"
	"github.com/Sanaruca/condominio/internal/core"
	"github.com/Sanaruca/condominio/internal/core/adapters/ozzo"
	"github.com/Sanaruca/condominio/internal/core/context"
	"github.com/Sanaruca/condominio/internal/core/usecase"
	"github.com/Sanaruca/condominio/internal/pagos/types/moneda"
	"github.com/Sanaruca/condominio/internal/services/tasa"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type GastoBase struct {
	Concepto string
	Monto    int
	Moneda   moneda.Moneda
	Fecha    *time.Time
	// TODO: comprobante
}

type RegistrarGastoDTO struct {
	GastoBase
	Proveedor string
}

type RegistrarGasto usecase.Handler[context.AdminContext, RegistrarGastoDTO, *gasto.Gasto]

type registrarGasto struct {
	repo        gasto.GastoRepository
	tasaService tasa.TasaService
}

func NewRegistrarGasto(repo gasto.GastoRepository, tasaService tasa.TasaService) RegistrarGasto {
	if repo == nil {
		panic("repo is nil")
	}
	if tasaService == nil {
		panic("tasaService is nil")
	}
	return registrarGasto{
		repo:        repo,
		tasaService: tasaService,
	}
}

func (uc registrarGasto) Exec(
	ctx context.AdminContext,
	input RegistrarGastoDTO,
) (*gasto.Gasto, core.Error) {

	// TODO: completar y testear

	if err := input.Validate(); err != nil {
		return nil, err
	}

	_tasa, err := uc.tasaService.ObtenerTasaParaFecha(tasa.CambioOficial, *input.Fecha)

	if err != nil {
		return nil, core.WrapError(err)
	}

	gasto, err := input.toGasto(ctx, _tasa.Valor)
	if err != nil {
		return nil, core.WrapError(err)
	}

	// TODO: buscar si ya existe un gasto con el mismo proveedor, monto, moneda y
	// fecha para evitar duplicados accidentales
	gastoID, err := uc.repo.Guardar(ctx, *gasto)

	if err != nil {
		return nil, core.WrapError(err)
	}

	return uc.repo.ObtenerPorID(ctx, *gastoID)
}

func (b *GastoBase) Validate() core.Error {

	if err := b.Moneda.Validate(); err != nil {
		return err
	}

	if b.Fecha == nil {
		now := time.Now()
		b.Fecha = &now
	}

	return ozzo.FirstOzzoErrorAdapter(b, validation.ValidateStruct(
		b,
		validation.Field(&b.Concepto, validation.Required, validation.Length(1, 100)),
		validation.Field(&b.Monto, validation.Required, validation.Min(1)),
	))
}

func (dto *RegistrarGastoDTO) Validate() core.Error {

	if err := dto.GastoBase.Validate(); err != nil {
		return err
	}

	err := validation.ValidateStruct(
		dto,
		validation.Field(&dto.Proveedor, validation.Required, validation.Length(1, 100)),
	)

	if err != nil {
		return ozzo.FirstOzzoErrorAdapter(dto, err)
	}

	return nil
}

func (dto RegistrarGastoDTO) toGasto(
	ctx context.AdminContext,
	tasa int,
) (*gasto.Gasto, core.Error) {
	return gasto.NuevoGasto(
		ctx.Session().Usuario().ID,
		dto.Concepto,
		dto.Proveedor,
		dto.Monto,
		dto.Moneda,
		tasa,
		*dto.Fecha,
	)

}
