// Deprecated: Este paquete pertenece al sistema legacy de pagos.
// Usar internal/transacciones/ en su lugar.
package pago

import (
	"time"

	"github.com/lucsky/cuid"

	"github.com/Sanaruca/condominio/internal/core/common/events"
	"github.com/Sanaruca/condominio/internal/core/common/filter"
	"github.com/Sanaruca/condominio/internal/core/common/quantity"
	"github.com/Sanaruca/condominio/internal/core/errors"
	"github.com/Sanaruca/condominio/internal/pagos/event"
	"github.com/Sanaruca/condominio/internal/pagos/types/metododepago"
	"github.com/Sanaruca/condominio/internal/pagos/types/moneda"
	"github.com/Sanaruca/condominio/internal/unidades/models/unidad"
)

var (
	ErrPagoNoEncontrado      = errors.New(errors.NOT_FOUND, "Pago no encontrado")
	ErrIdentificadorInvalido = errors.New(
		errors.INVALID_ARGUMENT,
		"El identificador no puede ser nulo",
	)
	ErrMontoInvalido     = errors.New(errors.INVALID_ARGUMENT, "El monto debe ser mayor a cero")
	ErrSaldoInsuficiente = errors.New(errors.INVALID_ARGUMENT, "monto excede el saldo del pago")
)

type Pago struct {
	id             string
	unidad         unidad.UnidadCodigo
	fecha          time.Time
	metodo         metododepago.MetodoDePago
	monto          quantity.Quantity
	referencia     *string
	moneda         moneda.Moneda
	tasa           quantity.Quantity
	firma          Firma
	destinos       []Destino
	event_notifier events.EventNotifier
}

func (p *Pago) ID() string                        { return p.id }
func (p *Pago) Unidad() unidad.UnidadCodigo       { return p.unidad }
func (p *Pago) Fecha() time.Time                  { return p.fecha }
func (p *Pago) Monto() quantity.Quantity          { return p.monto }
func (p *Pago) Moneda() moneda.Moneda             { return p.moneda }
func (p *Pago) Tasa() quantity.Quantity           { return p.tasa }
func (p *Pago) Referencia() *string               { return p.referencia }
func (p *Pago) Metodo() metododepago.MetodoDePago { return p.metodo }
func (p *Pago) Firma() Firma                      { return p.firma }
func (p *Pago) Destinos() []Destino               { return p.destinos }

// Deprecated: Use PagoFactory.Nuevo
func NuevoPago(
	unidad unidad.UnidadCodigo,
	fecha_de_pago time.Time,
	metodo metododepago.MetodoDePago,
	monto quantity.Quantity,
	moneda moneda.Moneda,
	tasa quantity.Quantity,
	referencia *string,
	usuarioID string,
) (*Pago, errors.CoreError) {

	firma, err := NuevaFirma(usuarioID)
	if err != nil {
		return nil, err
	}

	return newPago(
		cuid.New(),
		unidad,
		fecha_de_pago,
		metodo,
		monto,
		moneda,
		tasa,
		referencia,
		*firma,
	)

}

// Deprecated: Use PagoFactory.Assemble
func NuevoPagoFromStore(
	id string,
	unidad unidad.UnidadCodigo,
	fecha_de_pago time.Time,
	metodo metododepago.MetodoDePago,
	monto quantity.Quantity,
	moneda moneda.Moneda,
	tasa quantity.Quantity,
	referencia *string,
	firma Firma,
) (*Pago, errors.CoreError) {

	return newPago(
		id,
		unidad,
		fecha_de_pago,
		metodo,
		monto,
		moneda,
		tasa,
		referencia,
		firma,
	)
}

func (p *Pago) PullEvents() []events.Event { return p.event_notifier.Dispatch() }

// MontoDestinado retorna el total del monto del pago que ha sido destinado a deudas.
func (p *Pago) SaldoDestinado() quantity.Quantity {
	return p.Total().HappySub(p.SaldoDisponible())
}

// ObtenerDiferenciaDeDestinos retorna los destinos que tiene la entidad [Pago]
// pero que no están presentes en la lista de destinos_almacenados (IDs).
func (p *Pago) ObtenerDiferenciaDeDestinos(destinos_almacenados []string) []Destino {
	// Convertimos la lista a un Set (mapa) para búsquedas instantáneas
	existentes := make(map[string]struct{})
	for _, id := range destinos_almacenados {
		existentes[id] = struct{}{}
	}

	// Filtramos los destinos
	nuevos_destinos := make([]Destino, 0)
	for _, d := range p.destinos {
		if _, ok := existentes[d.ID()]; !ok {
			nuevos_destinos = append(nuevos_destinos, d)
		}
	}

	return nuevos_destinos
}

// AplicarPagoADeuda intenta cubrir el monto de una deuda utilizando el saldo disponible.
// Si el saldo es menor a la deuda, aplica el saldo restante como pago parcial.
// Retorna el monto efectivamente destinado a la deuda.
func (p *Pago) AplicarPagoADeuda(
	deudaID string,
	deuda_monto quantity.Quantity,
) (destinado quantity.Quantity, err errors.CoreError) {
	if p.EstaAgotado() {
		return destinado, ErrSaldoInsuficiente
	}

	saldo_disponible := p.SaldoDisponible()
	monto_a_destinar := deuda_monto

	// Si no alcanza para la deuda total, destinamos todo lo que queda
	if saldo_disponible.Value() < deuda_monto.Value() {
		monto_a_destinar = saldo_disponible
	}

	if err := p.AgregarDestino(deudaID, monto_a_destinar); err != nil {
		return destinado, err
	}

	return monto_a_destinar, nil

}

// AgregarDestino registra un movimiento de dinero hacia una deuda específica.
func (p *Pago) AgregarDestino(deudaID string, monto_a_destinar quantity.Quantity) errors.CoreError {
	if p.SaldoDisponible().Value() < monto_a_destinar.Value() {
		return ErrSaldoInsuficiente
	}

	destino, err := NuevoDestino(deudaID, monto_a_destinar)
	if err != nil {
		return err
	}

	p.destinos = append(p.destinos, *destino)

	return nil
}

func (p *Pago) EstaAgotado() bool {
	return p.SaldoDisponible().Value() <= 0
}

// SaldoDisponible calcula el saldo disponible del pago expresado en centavos
func (p *Pago) SaldoDisponible() quantity.Quantity {
	var aplicado quantity.Quantity
	for i, d := range p.destinos {

		if i == 0 {
			aplicado = d.destinado
			continue
		}

		aplicado = aplicado.HappyAdd(d.destinado)
	}

	return p.Total().HappySub(aplicado)
}

// Representa el monto total en la moneda base (USD)
// El monto está almacenado en centavos
func (p *Pago) Total() quantity.Quantity {
	switch p.moneda {
	case moneda.USD:
		return p.monto
	default:
		// Convertir de centavos en moneda local a centavos en USD
		return p.monto.HappyDiv(p.tasa)
	}
}

func newPago(
	id string,
	unidad unidad.UnidadCodigo,
	fecha_de_pago time.Time,
	metodo metododepago.MetodoDePago,
	monto quantity.Quantity,
	moneda moneda.Moneda,
	tasa quantity.Quantity,
	referencia *string,
	firma Firma,
) (*Pago, errors.CoreError) {

	if id == "" {
		return nil, ErrIdentificadorInvalido
	}

	if monto.Value() <= 0 {
		return nil, ErrMontoInvalido
	}

	if unidad == "" {
		return nil, errors.New(errors.INVALID_ARGUMENT, "La unidad no puede estar vacía")
	}

	if tasa.Value() <= 0 {
		return nil, errors.New(errors.INVALID_ARGUMENT, "La tasa debe ser mayor a cero")
	}

	if err := moneda.Validate(); err != nil {
		return nil, err
	}

	if err := firma.Validate(); err != nil {
		return nil, err
	}

	event_notifier := events.EventNotifier{}
	event_notifier.AddEvent(event.NewPagoRegistrado(id, monto, fecha_de_pago))

	return &Pago{
		id:             id,
		unidad:         unidad,
		fecha:          fecha_de_pago,
		metodo:         metodo,
		monto:          monto,
		referencia:     referencia,
		moneda:         moneda,
		tasa:           tasa,
		firma:          firma,
		event_notifier: event_notifier,
	}, nil

}

func (p Pago) FilterSpec() filter.Spec {
	return filter.Spec{
		"unidad": filter.TypeString,
	}
}
