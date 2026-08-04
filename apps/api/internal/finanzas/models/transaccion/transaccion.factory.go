package transaccion

import (
	"time"

	"github.com/lucsky/cuid"

	"github.com/Sanaruca/condominio/internal/core"
	"github.com/Sanaruca/condominio/internal/core/common/events"
	currency "github.com/Sanaruca/condominio/internal/core/common/moneda"
	"github.com/Sanaruca/condominio/internal/core/common/quantity"
	"github.com/Sanaruca/condominio/internal/core/errors"
	"github.com/Sanaruca/condominio/internal/finanzas/event"
	"github.com/Sanaruca/condominio/internal/finanzas/types/metodotransaccion"
	"github.com/Sanaruca/condominio/internal/finanzas/types/roldelmovimiento"
	"github.com/Sanaruca/condominio/internal/finanzas/types/tipodemovimiento"
)

type TransaccionFactory struct{}

func NewTransaccionFactory() *TransaccionFactory {
	return &TransaccionFactory{}
}

// Nuevo crea una transaccion generica validando que los movimientos cuadren.
func (f *TransaccionFactory) Nuevo(
	concepto string,
	monto_total quantity.Quantity,
	moneda currency.Moneda,
	metodo metodotransaccion.MetodoDeTransaccion,
	tasa quantity.Quantity,
	fecha time.Time,
	registrado_por string,
	cuotaID *string,
	movimientos []Movimiento,
) (*TransaccionFinanciera, core.Error) {
	if monto_total.Value() <= 0 {
		return nil, ErrMontoInvalido
	}
	if len(movimientos) == 0 {
		return nil, ErrMovimientosVacios
	}
	if err := f.validarAuditoria(monto_total, movimientos); err != nil {
		return nil, err
	}
	if err := moneda.Validate(); err != nil {
		return nil, err
	}
	if err := metodo.Validate(); err != nil {
		return nil, err
	}

	event_notifier := events.EventNotifier{}
	event_notifier.AddEvent(event.NewTransaccionRegistrada(
		cuid.New(),
		monto_total,
		fecha,
	))

	return &TransaccionFinanciera{
		id:             cuid.New(),
		fecha:          fecha,
		concepto:       concepto,
		monto_total:    monto_total,
		moneda:         moneda,
		metodo:         metodo,
		tasa:           tasa,
		registrado_por: registrado_por,
		registro:       time.Now().UTC(),
		movimientos:    movimientos,
		event_notifier: event_notifier,
	}, nil
}

// NuevoPago crea una transaccion de tipo pago (CREDITO hacia una UNIDAD).
func (f *TransaccionFactory) NuevoPago(
	unidad_codigo string,
	concepto string,
	monto_total quantity.Quantity,
	moneda currency.Moneda,
	metodo metodotransaccion.MetodoDeTransaccion,
	tasa quantity.Quantity,
	fecha time.Time,
	registrado_por string,
) (*TransaccionFinanciera, core.Error) {
	movimiento := newMovimiento(
		cuid.New(),
		tipodemovimiento.Credito,
		monto_total,
		roldelmovimiento.Unidad,
		&unidad_codigo,
		nil,
	)

	return f.Nuevo(
		concepto,
		monto_total,
		moneda,
		metodo,
		tasa,
		fecha,
		registrado_por,
		nil,
		[]Movimiento{*movimiento},
	)
}

// NuevoGasto crea una transaccion de tipo gasto (DEBITO hacia un PROVEEDOR o CONDOMINIO).
func (f *TransaccionFactory) NuevoGasto(
	concepto string,
	proveedor_id *string,
	es_condominio bool,
	monto_total quantity.Quantity,
	moneda currency.Moneda,
	metodo metodotransaccion.MetodoDeTransaccion,
	tasa quantity.Quantity,
	fecha time.Time,
	cuota_id *string,
	registrado_por string,
) (*TransaccionFinanciera, core.Error) {
	if cuota_id != nil && *cuota_id == "" {
		cuota_id = nil
	}

	var rol roldelmovimiento.RolDelMovimiento
	if es_condominio {
		rol = roldelmovimiento.Condominio
	} else {
		rol = roldelmovimiento.Proveedor
	}

	movimiento := newMovimiento(
		cuid.New(),
		tipodemovimiento.Debito,
		monto_total,
		rol,
		nil,
		proveedor_id,
	)

	return f.Nuevo(
		concepto,
		monto_total,
		moneda,
		metodo,
		tasa,
		fecha,
		registrado_por,
		cuota_id,
		[]Movimiento{*movimiento},
	)
}

// NuevoReembolso crea una transaccion de reembolso (DEBITO hacia una UNIDAD).
func (f *TransaccionFactory) NuevoReembolso(
	unidad_codigo string,
	concepto string,
	monto_total quantity.Quantity,
	moneda currency.Moneda,
	metodo metodotransaccion.MetodoDeTransaccion,
	tasa quantity.Quantity,
	fecha time.Time,
	registrado_por string,
) (*TransaccionFinanciera, core.Error) {
	movimiento := newMovimiento(
		cuid.New(),
		tipodemovimiento.Debito,
		monto_total,
		roldelmovimiento.Unidad,
		&unidad_codigo,
		nil,
	)

	return f.Nuevo(
		concepto,
		monto_total,
		moneda,
		metodo,
		tasa,
		fecha,
		registrado_por,
		nil,
		[]Movimiento{*movimiento},
	)
}

// Assemble reconstruye una transaccion desde la persistencia (infalible).
func (f *TransaccionFactory) Assemble(
	id string,
	fecha time.Time,
	concepto string,
	monto_total quantity.Quantity,
	moneda currency.Moneda,
	metodo metodotransaccion.MetodoDeTransaccion,
	tasa quantity.Quantity,
	registrado_por string,
	registro time.Time,
	movimientos []Movimiento,
) *TransaccionFinanciera {
	return &TransaccionFinanciera{
		id:             id,
		fecha:          fecha,
		concepto:       concepto,
		monto_total:    monto_total,
		moneda:         moneda,
		metodo:         metodo,
		tasa:           tasa,
		registrado_por: registrado_por,
		registro:       registro,
		movimientos:    movimientos,
	}
}

func (f *TransaccionFactory) validarAuditoria(
	monto_total quantity.Quantity,
	movimientos []Movimiento,
) core.Error {
	var suma quantity.Quantity
	for i, m := range movimientos {
		if i == 0 {
			suma = m.monto
		} else {
			suma = suma.HappyAdd(m.monto)
		}
	}

	if suma.Value() != monto_total.Value() {
		return ErrDescuadreContable
	}

	return nil
}

// Errors propios del dominio
var (
	ErrIdentificadorInvalido = core.NewError(
		errors.INVALID_ARGUMENT,
		"El identificador no puede ser nulo",
	)
	ErrTasaRequerida = core.NewError(
		errors.INVALID_ARGUMENT,
		"La tasa debe ser mayor a cero cuando la moneda es VED",
	)
)
