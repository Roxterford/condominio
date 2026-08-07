package transaccion

import (
	"time"

	"github.com/Sanaruca/condominio/internal/core/common/events"
	"github.com/Sanaruca/condominio/internal/core/common/filter"
	currency "github.com/Sanaruca/condominio/internal/core/common/moneda"
	"github.com/Sanaruca/condominio/internal/core/common/quantity"
	"github.com/Sanaruca/condominio/internal/core/errors"
	"github.com/Sanaruca/condominio/internal/finanzas/types/metodotransaccion"
)

var (
	ErrTransaccionNoEncontrada = errors.New(errors.NOT_FOUND, "Transaccion no encontrada")
	ErrMontoInvalido           = errors.New(
		errors.INVALID_ARGUMENT,
		"El monto total debe ser mayor a cero",
	)
	ErrMovimientosVacios = errors.New(
		errors.INVALID_ARGUMENT,
		"La transaccion debe tener al menos un movimiento",
	)
	ErrDescuadreContable = errors.New(
		errors.INVALID_ARGUMENT,
		"La suma de los movimientos no coincide con el monto total de la transaccion",
	)
	ErrCuotaRequerida = errors.New(
		errors.INVALID_ARGUMENT,
		"El ID de cuota es requerido para gastos del condominio",
	)
)

type TransaccionFinanciera struct {
	id             string
	fecha          time.Time
	concepto       string
	monto_total    quantity.Quantity
	moneda         currency.Moneda
	metodo         metodotransaccion.MetodoDeTransaccion
	tasa           quantity.Quantity
	registrado_por string
	registro       time.Time
	movimientos    []Movimiento
	event_notifier events.EventNotifier
}

func (t *TransaccionFinanciera) ID() string       { return t.id }
func (t *TransaccionFinanciera) Fecha() time.Time { return t.fecha }
func (t *TransaccionFinanciera) Concepto() string { return t.concepto }

func (t *TransaccionFinanciera) MontoTotal() quantity.Quantity                 { return t.monto_total }
func (t *TransaccionFinanciera) Moneda() currency.Moneda                       { return t.moneda }
func (t *TransaccionFinanciera) Metodo() metodotransaccion.MetodoDeTransaccion { return t.metodo }
func (t *TransaccionFinanciera) Tasa() quantity.Quantity                       { return t.tasa }

func (t *TransaccionFinanciera) RegistradoPor() string { return t.registrado_por }
func (t *TransaccionFinanciera) Registro() time.Time   { return t.registro }

func (t *TransaccionFinanciera) Movimientos() []Movimiento { return t.movimientos }

// TotalEnUSD retorna el monto total convertido a USD
func (t *TransaccionFinanciera) TotalEnUSD() quantity.Quantity {
	switch t.moneda {
	case currency.USD:
		return t.monto_total
	default:
		return t.monto_total.HappyDiv(t.tasa)
	}
}

func (t *TransaccionFinanciera) PullEvents() []events.Event {
	return t.event_notifier.Dispatch()
}

func (t TransaccionFinanciera) FilterSpec() filter.Spec {
	return filter.Spec{
		"concepto": filter.TypeString,
		"cuota":    filter.TypeString,
		"moneda":   filter.TypeString,
	}
}
