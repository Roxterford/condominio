package pago

import (
	"time"

	"github.com/lucsky/cuid"

	"github.com/Sanaruca/condominio/internal/core/common/events"
	"github.com/Sanaruca/condominio/internal/core/common/quantity"
	"github.com/Sanaruca/condominio/internal/core/errors"
	"github.com/Sanaruca/condominio/internal/pagos/event"
	"github.com/Sanaruca/condominio/internal/pagos/types/metododepago"
	"github.com/Sanaruca/condominio/internal/pagos/types/moneda"
	"github.com/Sanaruca/condominio/internal/unidades/models/unidad"
)

type PagoFactory struct {
}

func NewPagoFactory() *PagoFactory {
	return &PagoFactory{}
}

// Nuevo crea una instancia válida de Pago.
func (*PagoFactory) Nuevo(
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

func (*PagoFactory) Assemble(
	id string,
	unidad unidad.UnidadCodigo,
	fecha_de_pago time.Time,
	metodo metododepago.MetodoDePago,
	monto quantity.Quantity,
	moneda moneda.Moneda,
	tasa quantity.Quantity,
	referencia *string,
	registro time.Time,
	registrado_por string,
	actualizacion time.Time,
	actualizado_por string,
) *Pago {

	event_notifier := events.EventNotifier{}
	event_notifier.AddEvent(event.NewPagoRegistrado(id, monto, fecha_de_pago))

	firma, _ := NuevaFirmaFromStore(registro, registrado_por, actualizacion, actualizado_por)

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
	}

}
