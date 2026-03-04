package pagos

import (
	"time"

	"github.com/Sanaruca/condominio/internal/core/common/events"
	"github.com/Sanaruca/condominio/internal/core/errors"
	"github.com/Sanaruca/condominio/internal/pagos/event"
	"github.com/Sanaruca/condominio/internal/pagos/types/metododepago"
	"github.com/Sanaruca/condominio/internal/pagos/types/moneda"
	"github.com/lucsky/cuid"
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
	villa          int
	fecha          time.Time
	metodo         metododepago.MetodoDePago
	monto          int
	referencia     *string
	moneda         moneda.Moneda
	tasa           int
	firma          Firma
	destinos       []Destino
	event_notifier events.EventNotifier
}

func (p *Pago) ID() string                        { return p.id }
func (p *Pago) Villa() int                        { return p.villa }
func (p *Pago) Fecha() time.Time                  { return p.fecha }
func (p *Pago) Monto() int                        { return p.monto }
func (p *Pago) Moneda() moneda.Moneda             { return p.moneda }
func (p *Pago) Tasa() int                         { return p.tasa }
func (p *Pago) Referencia() *string               { return p.referencia }
func (p *Pago) Metodo() metododepago.MetodoDePago { return p.metodo }
func (p *Pago) Firma() Firma                      { return p.firma }
func (p *Pago) Destinos() []Destino               { return p.destinos }

// NuevoPago crea una instancia válida de Pago.
func NuevoPago(
	villa int,
	fecha_de_pago time.Time,
	metodo metododepago.MetodoDePago,
	monto int,
	moneda moneda.Moneda,
	tasa int,
	referencia *string,
	usuarioID string,
) (*Pago, errors.CoreError) {

	firma, err := NuevaFirma(usuarioID)
	if err != nil {
		return nil, err
	}

	return newPago(
		cuid.New(),
		villa,
		fecha_de_pago,
		metodo,
		monto,
		moneda,
		tasa,
		referencia,
		*firma,
	)

}

func NuevoPagoFromStore(
	id string,
	villa int,
	fecha_de_pago time.Time,
	metodo metododepago.MetodoDePago,
	monto int,
	moneda moneda.Moneda,
	tasa int,
	referencia *string,
	firma Firma,
) (*Pago, errors.CoreError) {

	return newPago(
		id,
		villa,
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
	monto_deuda int,
) (destinado int, err errors.CoreError) {
	if p.EstaAgotado() {
		return 0, ErrSaldoInsuficiente
	}

	saldo_disponible := p.SaldoDisponible()
	monto_a_destinar := monto_deuda

	// Si no alcanza para la deuda total, destinamos todo lo que queda
	if saldo_disponible < monto_deuda {
		monto_a_destinar = saldo_disponible
	}

	if err := p.AgregarDestino(deudaID, monto_a_destinar); err != nil {
		return 0, err
	}

	return monto_a_destinar, nil

}

// AgregarDestino registra un movimiento de dinero hacia una deuda específica.
func (p *Pago) AgregarDestino(deudaID string, monto_a_destinar int) errors.CoreError {
	if p.SaldoDisponible() < monto_a_destinar {
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
	return p.SaldoDisponible() < 1
}

// SaldoDisponible calcula el saldo disponible del pago expresado en centavos
func (p *Pago) SaldoDisponible() int {
	aplicado := 0
	for _, d := range p.destinos {
		aplicado += d.Destinado()
	}
	return p.monto - aplicado
}

// Representa el monto total en la moneda base (USD)
// El monto está almacenado en centavos
func (p *Pago) Total() int {
	switch p.moneda {
	case moneda.USD:
		return p.monto
	default:
		// Convertir de centavos en moneda local a centavos en USD
		// Ejemplo: 10000 centavos locales / 1000 tasa = 10 USD = 1000 centavos USD
		return int(float64(p.monto) / float64(p.tasa) * 100)
	}
}

func newPago(
	id string,
	villa int,
	fecha_de_pago time.Time,
	metodo metododepago.MetodoDePago,
	monto int,
	moneda moneda.Moneda,
	tasa int,
	referencia *string,
	firma Firma,
) (*Pago, errors.CoreError) {

	if id == "" {
		return nil, ErrIdentificadorInvalido
	}

	if monto <= 0 {
		return nil, ErrMontoInvalido
	}

	if villa <= 0 {
		return nil, errors.New(errors.INVALID_ARGUMENT, "La villa debe ser mayor a cero")
	}

	if tasa <= 0 {
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
		villa:          villa,
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
