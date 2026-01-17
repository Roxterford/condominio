package service

import (
	"github.com/Sanaruca/condominio/internal/pagos/app"
	"github.com/Sanaruca/condominio/internal/pagos/app/command"
)

type PagoService struct {
	Commands app.Commands
}

func New() *PagoService {

	registrarPago := command.NewRegistrarPago()
	registrarPagoADeuda := command.NewRegistarPagoADeuda()

	return &PagoService{
		Commands: app.Commands{
			RegistrarPago:       registrarPago,
			RegistrarPagoADeuda: registrarPagoADeuda,
		},
	}
}
