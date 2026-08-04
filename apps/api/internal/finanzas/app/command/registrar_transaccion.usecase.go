package command

import (
	"time"

	"github.com/Sanaruca/condominio/internal/core"
	"github.com/Sanaruca/condominio/internal/core/common/moneda"
	"github.com/Sanaruca/condominio/internal/core/common/quantity"
	cc "github.com/Sanaruca/condominio/internal/core/context"
	"github.com/Sanaruca/condominio/internal/core/usecase"
	"github.com/Sanaruca/condominio/internal/finanzas/models/transaccion"
	"github.com/Sanaruca/condominio/internal/finanzas/types/metodotransaccion"
	"github.com/Sanaruca/condominio/internal/unidades/models/unidad"
)

type TipoTransaccion string

const (
	TipoPago      TipoTransaccion = "pago"
	TipoGasto     TipoTransaccion = "gasto"
	TipoReembolso TipoTransaccion = "reembolso"
)

type RegistrarTransaccionDTO struct {
	Tipo         TipoTransaccion
	Unidad       *unidad.UnidadCodigo
	Proveedor    *string
	Concepto     string
	Monto        int
	Moneda       moneda.Moneda
	Metodo       metodotransaccion.MetodoDeTransaccion
	Tasa         int
	Fecha        *time.Time
	Referencia   *string
	CuotaID      *string
	EsCondominio bool
}

type RegistrarTransaccion usecase.WithContextInput[cc.AdminContext, RegistrarTransaccionDTO]

type registrarTransaccion struct {
	repo     transaccion.TransaccionRepository
	factory  *transaccion.TransaccionFactory
	unidades unidad.UnidadRepository
	qf       *quantity.QuantityFactory
}

func NewRegistrarTransaccion(
	repo transaccion.TransaccionRepository,
	factory *transaccion.TransaccionFactory,
	unidadRepo unidad.UnidadRepository,
	quantityFactory *quantity.QuantityFactory,
) RegistrarTransaccion {
	if repo == nil {
		panic("repo is nil")
	}
	if factory == nil {
		panic("factory is nil")
	}
	if unidadRepo == nil {
		panic("unidadRepo is nil")
	}
	if quantityFactory == nil {
		panic("qf is nil")
	}
	return &registrarTransaccion{
		repo:     repo,
		factory:  factory,
		unidades: unidadRepo,
		qf:       quantityFactory,
	}
}

func (uc *registrarTransaccion) Exec(
	ctx cc.AdminContext,
	input RegistrarTransaccionDTO,
) (any, core.Error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	fecha := *input.Fecha

	tasaVal := int64(input.Tasa)
	monto := uc.qf.Assemble(int64(input.Monto))
	tasaQ := uc.qf.Assemble(tasaVal)

	var t *transaccion.TransaccionFinanciera
	var err error

	switch input.Tipo {
	case TipoPago:
		if input.Unidad == nil {
			return nil, core.NewValidationError("La unidad es requerida para pagos")
		}
		unidadCodigo := string(*input.Unidad)
		exists, e := uc.unidades.ExistsCodigo(ctx, *input.Unidad)
		if e != nil {
			return nil, e
		}
		if !exists {
			return nil, unidad.ErrUnidadNoEncontrada
		}
		t, err = uc.factory.NuevoPago(
			unidadCodigo,
			input.Concepto,
			monto,
			input.Moneda,
			input.Metodo,
			tasaQ,
			fecha,
			ctx.Session().Usuario().ID,
		)

	case TipoGasto:
		t, err = uc.factory.NuevoGasto(
			input.Concepto,
			input.Proveedor,
			input.EsCondominio,
			monto,
			input.Moneda,
			input.Metodo,
			tasaQ,
			fecha,
			input.CuotaID,
			ctx.Session().Usuario().ID,
		)

	case TipoReembolso:
		if input.Unidad == nil {
			return nil, core.NewValidationError("La unidad es requerida para reembolsos")
		}
		t, err = uc.factory.NuevoReembolso(
			string(*input.Unidad),
			input.Concepto,
			monto,
			input.Moneda,
			input.Metodo,
			tasaQ,
			fecha,
			ctx.Session().Usuario().ID,
		)

	default:
		return nil, core.NewValidationError("Tipo de transaccion no valido")
	}

	if err != nil {
		return nil, core.WrapError(err)
	}

	if err := uc.repo.Guardar(ctx, t); err != nil {
		return nil, err
	}

	return t, nil
}

func (dto *RegistrarTransaccionDTO) Validate() core.Error {

	if err := dto.Metodo.Validate(); err != nil {
		return err
	}

	if err := dto.Moneda.Validate(); err != nil {
		return err
	}

	if dto.Moneda == moneda.VED && dto.Tasa < 1 {
		return core.NewValidationError("La tasa es requerida para transacciones en VED")
	}

	if dto.Fecha == nil {
		now := time.Now()
		dto.Fecha = &now
	}

	if dto.Fecha.After(time.Now()) {
		return core.NewInvalidArgumentError("La fecha no puede ser futura")
	}

	if dto.Concepto == "" {
		return core.NewValidationError("El concepto es requerido")
	}
	if dto.Monto < 1 {
		return core.NewValidationError("El monto debe ser mayor a 0")
	}

	return nil
}
