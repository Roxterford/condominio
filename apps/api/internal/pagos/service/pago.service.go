package service

import (
	"github.com/Sanaruca/condominio/internal/core/common/events"
	"github.com/Sanaruca/condominio/internal/core/common/quantity"
	"github.com/Sanaruca/condominio/internal/pagos/app"
	"github.com/Sanaruca/condominio/internal/pagos/app/command"
	"github.com/Sanaruca/condominio/internal/pagos/models/pago"
	"github.com/Sanaruca/condominio/internal/services/tasa"
	"github.com/Sanaruca/condominio/internal/unidades/models/unidad"
)

type PagoService struct {
	Commands app.Commands
}

func New(
	pago_repository pago.PagoRepository,
	unidad_repository unidad.UnidadRepository,
	event_bus events.EventBus,
	tasa_service tasa.TasaService,
	quantity_factory *quantity.QuantityFactory,
) *PagoService {

	registrarPago := command.NewRegistrarPago(
		pago_repository,
		unidad_repository,
		event_bus,
		tasa_service,
		quantity_factory,
	)

	return &PagoService{
		Commands: app.Commands{
			RegistrarPago: registrarPago,
		},
	}
}
