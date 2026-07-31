// Deprecated: Servicio legacy de pagos. Usar internal/transacciones/service/ en su lugar.
package service

import (
	"github.com/Sanaruca/condominio/internal/core/common/events"
	"github.com/Sanaruca/condominio/internal/core/common/quantity"
	"github.com/Sanaruca/condominio/internal/pagos/app"
	"github.com/Sanaruca/condominio/internal/pagos/app/command"
	"github.com/Sanaruca/condominio/internal/pagos/app/query"
	"github.com/Sanaruca/condominio/internal/pagos/models/pago"
	"github.com/Sanaruca/condominio/internal/services/tasa"
	"github.com/Sanaruca/condominio/internal/unidades/models/unidad"
)

type PagoService struct {
	Commands app.Commands
	Queries  app.Queries
}

func New(
	pago_repository pago.PagoRepository,
	pago_factory *pago.PagoFactory,
	unidad_repository unidad.UnidadRepository,
	event_bus events.EventBus,
	tasa_service tasa.TasaService,
	quantity_factory *quantity.QuantityFactory,
	aplicarPago command.AplicarPago,
) *PagoService {

	registrarPago := command.NewRegistrarPago(
		pago_repository,
		pago_factory,
		unidad_repository,
		event_bus,
		tasa_service,
		quantity_factory,
	)

	reprocesarPagosHuerfanos := command.NewReprocesarPagosHuerfanos(
		pago_repository,
		aplicarPago,
	)

	return &PagoService{
		Commands: app.Commands{
			RegistrarPago:            registrarPago,
			ReprocesarPagosHuerfanos: reprocesarPagosHuerfanos,
		},
		Queries: app.Queries{
			ObtenerPagos: query.NewObtenerPagos(pago_repository),
		},
	}
}
