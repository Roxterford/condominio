package operacion

import (
	"time"

	"github.com/lucsky/cuid"

	"github.com/Sanaruca/condominio/internal/core"
	"github.com/Sanaruca/condominio/internal/core/common/events"
	currency "github.com/Sanaruca/condominio/internal/core/common/moneda"
	"github.com/Sanaruca/condominio/internal/core/common/quantity"
	"github.com/Sanaruca/condominio/internal/finanzas/event"
	"github.com/Sanaruca/condominio/internal/finanzas/types/metodoperacion"
	"github.com/Sanaruca/condominio/internal/finanzas/types/roldestionoperacion"
	"github.com/Sanaruca/condominio/internal/finanzas/types/tipoperacion"
)

type OperacionFactory struct{}

func NewOperacionFactory() *OperacionFactory {
	return &OperacionFactory{}
}

// Nuevo es el guardian de la integridad: valida los datos de negocio vigentes
// y genera un identificador y registro actuales. Se usa al introducir datos nuevos.
func (f *OperacionFactory) Nuevo(
	concepto string,
	monto quantity.Quantity,
	moneda currency.Moneda,
	metodo metodoperacion.MetodoDeOperacion,
	tasa quantity.Quantity,
	tipo tipoperacion.TipoDeOperacion,
	rol roldestionoperacion.RolDestinoDeOperacion,
	cuota *string,
	unidad_codigo *string,
	proveedor_id *string,
	fecha time.Time,
	registrado_por string,
) (*Operacion, core.Error) {
	if monto.Value() <= 0 {
		return nil, ErrMontoInvalido
	}
	if err := moneda.Validate(); err != nil {
		return nil, err
	}
	if err := metodo.Validate(); err != nil {
		return nil, err
	}
	if err := tipo.Validate(); err != nil {
		return nil, err
	}
	if err := f.validarRol(rol, unidad_codigo, proveedor_id); err != nil {
		return nil, err
	}

	notifier := events.EventNotifier{}
	notifier.AddEvent(event.NewOperacionRegistrada(cuid.New(), monto, fecha))

	return &Operacion{
		id:             cuid.New(),
		fecha:          fecha,
		concepto:       concepto,
		monto:          monto,
		moneda:         moneda,
		metodo:         metodo,
		tasa:           tasa,
		tipo:           tipo,
		rol:            rol,
		cuota:          cuota,
		unidad_codigo:  unidad_codigo,
		proveedor_id:   proveedor_id,
		registrado_por: registrado_por,
		registro:       time.Now().UTC(),
		event_notifier: notifier,
	}, nil
}

// NuevoPago registra un pago (CREDITO hacia una UNIDAD): dinero que entra.
func (f *OperacionFactory) NuevoPago(
	unidad_codigo string,
	concepto string,
	monto quantity.Quantity,
	moneda currency.Moneda,
	metodo metodoperacion.MetodoDeOperacion,
	tasa quantity.Quantity,
	fecha time.Time,
	registrado_por string,
) (*Operacion, core.Error) {
	return f.Nuevo(
		concepto,
		monto,
		moneda,
		metodo,
		tasa,
		tipoperacion.Credito,
		roldestionoperacion.Unidad,
		nil,
		&unidad_codigo,
		nil,
		fecha,
		registrado_por,
	)
}

// NuevoGasto registra un gasto (DEBITO hacia un PROVEEDOR o CONDOMINIO): dinero que sale.
func (f *OperacionFactory) NuevoGasto(
	concepto string,
	proveedor_id *string,
	es_condominio bool,
	monto quantity.Quantity,
	moneda currency.Moneda,
	metodo metodoperacion.MetodoDeOperacion,
	tasa quantity.Quantity,
	fecha time.Time,
	cuota_id *string,
	registrado_por string,
) (*Operacion, core.Error) {
	if cuota_id != nil && *cuota_id == "" {
		cuota_id = nil
	}

	var rol roldestionoperacion.RolDestinoDeOperacion
	var unidad_codigo *string

	if es_condominio {
		rol = roldestionoperacion.Condominio
	} else {
		if proveedor_id == nil || *proveedor_id == "" {
			return nil, ErrProveedorRequerido
		}
		rol = roldestionoperacion.Proveedor
	}

	return f.Nuevo(
		concepto,
		monto,
		moneda,
		metodo,
		tasa,
		tipoperacion.Debito,
		rol,
		cuota_id,
		unidad_codigo,
		proveedor_id,
		fecha,
		registrado_por,
	)
}

// NuevoReembolso registra un reembolso (DEBITO hacia una UNIDAD): dinero que sale
// de vuelta a una unidad.
func (f *OperacionFactory) NuevoReembolso(
	unidad_codigo string,
	concepto string,
	monto quantity.Quantity,
	moneda currency.Moneda,
	metodo metodoperacion.MetodoDeOperacion,
	tasa quantity.Quantity,
	fecha time.Time,
	registrado_por string,
) (*Operacion, core.Error) {
	return f.Nuevo(
		concepto,
		monto,
		moneda,
		metodo,
		tasa,
		tipoperacion.Debito,
		roldestionoperacion.Unidad,
		nil,
		&unidad_codigo,
		nil,
		fecha,
		registrado_por,
	)
}

// Assemble reconstruye una operacion desde la persistencia (infalible).
func (f *OperacionFactory) Assemble(
	id string,
	fecha time.Time,
	concepto string,
	monto quantity.Quantity,
	moneda currency.Moneda,
	metodo metodoperacion.MetodoDeOperacion,
	tasa quantity.Quantity,
	tipo tipoperacion.TipoDeOperacion,
	rol roldestionoperacion.RolDestinoDeOperacion,
	cuota *string,
	unidad_codigo *string,
	proveedor_id *string,
	registrado_por string,
	registro time.Time,
) *Operacion {
	return &Operacion{
		id:             id,
		fecha:          fecha,
		concepto:       concepto,
		monto:          monto,
		moneda:         moneda,
		metodo:         metodo,
		tasa:           tasa,
		tipo:           tipo,
		rol:            rol,
		cuota:          cuota,
		unidad_codigo:  unidad_codigo,
		proveedor_id:   proveedor_id,
		registrado_por: registrado_por,
		registro:       registro,
	}
}

func (f *OperacionFactory) validarRol(
	rol roldestionoperacion.RolDestinoDeOperacion,
	unidad_codigo *string,
	proveedor_id *string,
) core.Error {
	switch rol {
	case roldestionoperacion.Unidad:
		if unidad_codigo == nil || *unidad_codigo == "" {
			return ErrUnidadRequerida
		}
	case roldestionoperacion.Proveedor:
		if proveedor_id == nil || *proveedor_id == "" {
			return ErrProveedorRequerido
		}
	case roldestionoperacion.Condominio:
	default:
		return ErrRolInvalido
	}

	return nil
}
