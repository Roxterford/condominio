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
	"github.com/Sanaruca/condominio/internal/finanzas/types/roldestinoperacion"
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
	rol roldestinoperacion.RolDestinoDeOperacion,
	cuota *string,
	unidad_codigo *string,
	proveedor_id *string,
	fecha time.Time,
	registrado_por string,
	correlationID string,
	causationID string,
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

	notifier := events.PendingEvents{}
	notifier.AddEvent(
		event.NewOperacionRegistrada(cuid.New(), monto, fecha, correlationID, causationID),
	)

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
	correlationID string,
	causationID string,
) (*Operacion, core.Error) {
	return f.Nuevo(
		concepto,
		monto,
		moneda,
		metodo,
		tasa,
		tipoperacion.Credito,
		roldestinoperacion.Unidad,
		nil,
		&unidad_codigo,
		nil,
		fecha,
		registrado_por,
		correlationID,
		causationID,
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
	correlationID string,
	causationID string,
) (*Operacion, core.Error) {
	if cuota_id != nil && *cuota_id == "" {
		cuota_id = nil
	}

	var rol roldestinoperacion.RolDestinoDeOperacion
	var unidad_codigo *string

	if es_condominio {
		rol = roldestinoperacion.Condominio
	} else {
		if proveedor_id == nil || *proveedor_id == "" {
			return nil, ErrProveedorRequerido
		}
		rol = roldestinoperacion.Proveedor
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
		correlationID,
		causationID,
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
	correlationID string,
	causationID string,
) (*Operacion, core.Error) {
	return f.Nuevo(
		concepto,
		monto,
		moneda,
		metodo,
		tasa,
		tipoperacion.Debito,
		roldestinoperacion.Unidad,
		nil,
		&unidad_codigo,
		nil,
		fecha,
		registrado_por,
		correlationID,
		causationID,
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
	rol roldestinoperacion.RolDestinoDeOperacion,
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
		event_notifier: events.PendingEvents{},
	}
}

func (f *OperacionFactory) validarRol(
	rol roldestinoperacion.RolDestinoDeOperacion,
	unidad_codigo *string,
	proveedor_id *string,
) core.Error {
	switch rol {
	case roldestinoperacion.Unidad:
		if unidad_codigo == nil || *unidad_codigo == "" {
			return ErrUnidadRequerida
		}
	case roldestinoperacion.Proveedor:
		if proveedor_id == nil || *proveedor_id == "" {
			return ErrProveedorRequerido
		}
	case roldestinoperacion.Condominio:
	default:
		return ErrRolInvalido
	}

	return nil
}
