package service

import (
	"github.com/Sanaruca/condominio/internal/core/common/events"
	"github.com/Sanaruca/condominio/internal/pagos"
	"github.com/Sanaruca/condominio/internal/pagos/app"
	"github.com/Sanaruca/condominio/internal/pagos/app/command"
	"github.com/Sanaruca/condominio/internal/services/tasa"
	"github.com/Sanaruca/condominio/internal/villas"
)

type PagoService struct {
	Commands app.Commands
}

func New(
	pago_repository pagos.PagoRepository,
	villa_repository villas.VillaRepository,
	event_bus events.EventBus,
	tasa_service tasa.TasaService,
) *PagoService {

	registrarPago := command.NewRegistrarPago(
		pago_repository,
		villa_repository,
		event_bus,
		tasa_service,
	)

	return &PagoService{
		Commands: app.Commands{
			RegistrarPago: registrarPago,
		},
	}
}
