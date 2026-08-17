package operacion

import (
	"time"

	"github.com/Sanaruca/condominio/internal/core"
	"github.com/Sanaruca/condominio/internal/core/common/events"
	"github.com/Sanaruca/condominio/internal/core/common/filter"
	"github.com/Sanaruca/condominio/internal/core/common/moneda"
	currency "github.com/Sanaruca/condominio/internal/core/common/moneda"
	"github.com/Sanaruca/condominio/internal/core/common/quantity"
	"github.com/Sanaruca/condominio/internal/core/errors"
	"github.com/Sanaruca/condominio/internal/finanzas/types/metodoperacion"
	"github.com/Sanaruca/condominio/internal/finanzas/types/roldestionoperacion"
	"github.com/Sanaruca/condominio/internal/finanzas/types/tipoperacion"
)

var (
	ErrOperacionNoEncontrada = errors.New(errors.NOT_FOUND, "Operacion no encontrada")
	ErrMontoInvalido         = errors.New(
		errors.INVALID_ARGUMENT,
		"El monto debe ser mayor a cero",
	)
	ErrRolInvalido = errors.New(
		errors.INVALID_ARGUMENT,
		"El rol de la operacion es invalido",
	)
	ErrUnidadRequerida = errors.New(
		errors.INVALID_ARGUMENT,
		"El codigo de unidad es requerido cuando el rol es UNIDAD",
	)
	ErrProveedorRequerido = errors.New(
		errors.INVALID_ARGUMENT,
		"El proveedor es requerido cuando el rol es PROVEEDOR",
	)
)

// Operacion es la entidad de primera clase de finanzas. Representa un registro
// individual (pago, gasto, reembolso, entrada o salida) con su metadata completa.
// Resultado de la fusion del antiguo Movimiento con su Transaccion.
type Operacion struct {
	id             string
	fecha          time.Time
	concepto       string
	monto          quantity.Quantity
	moneda         currency.Moneda
	metodo         metodoperacion.MetodoDeOperacion
	tasa           quantity.Quantity
	tipo           tipoperacion.TipoDeOperacion
	rol            roldestionoperacion.RolDestinoDeOperacion
	cuota          *string
	unidad_codigo  *string
	proveedor_id   *string
	registrado_por string
	registro       time.Time
	event_notifier events.EventNotifier
}

func (o *Operacion) ID() string                               { return o.id }
func (o *Operacion) Fecha() time.Time                         { return o.fecha }
func (o *Operacion) Concepto() string                         { return o.concepto }
func (o *Operacion) Monto() quantity.Quantity                 { return o.monto }
func (o *Operacion) Moneda() currency.Moneda                  { return o.moneda }
func (o *Operacion) Metodo() metodoperacion.MetodoDeOperacion { return o.metodo }
func (o *Operacion) Tasa() quantity.Quantity                  { return o.tasa }
func (o *Operacion) Total() quantity.Quantity {
	if o.moneda == moneda.USD {
		return o.monto
	}
	return o.monto.HappyDiv(o.tasa)
}
func (o *Operacion) Tipo() tipoperacion.TipoDeOperacion             { return o.tipo }
func (o *Operacion) Rol() roldestionoperacion.RolDestinoDeOperacion { return o.rol }
func (o *Operacion) CuotaID() *string                               { return o.cuota }
func (o *Operacion) UnidadCodigo() *string                          { return o.unidad_codigo }
func (o *Operacion) ProveedorID() *string                           { return o.proveedor_id }
func (o *Operacion) RegistradoPor() string                          { return o.registrado_por }
func (o *Operacion) Registro() time.Time                            { return o.registro }

func (o *Operacion) EsDebito() bool  { return o.tipo.DEBITO() }
func (o *Operacion) EsCredito() bool { return o.tipo.CREDITO() }

// MontoEnUSD retorna el monto convertido a USD segun la moneda y tasa de la operacion.
func (o *Operacion) MontoEnUSD() quantity.Quantity {
	switch o.moneda {
	case currency.USD:
		return o.monto
	default:
		return o.monto.HappyDiv(o.tasa)
	}
}

func (o *Operacion) AsignarCuota(cuotaID string) core.Error {

	if o.cuota != nil {
		return core.NewValidationError("Operacion ya posee cuota asignada")
	}

	if cuotaID == "" {
		return core.NewValidationError("Asignacion de cuota no es valida")
	}

	o.cuota = &cuotaID

	return nil
}

func (o *Operacion) PullEvents() []events.Event {
	return o.event_notifier.Dispatch()
}

func (o Operacion) FilterSpec() filter.Spec {
	return filter.Spec{
		"concepto":      filter.TypeString,
		"tipo":          filter.TypeString,
		"rol":           filter.TypeString,
		"cuota":         filter.TypeString,
		"moneda":        filter.TypeString,
		"unidad_codigo": filter.TypeString,
		"proveedor_id":  filter.TypeString,
	}
}
