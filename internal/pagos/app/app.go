package app

import "github.com/Sanaruca/condominio/internal/pagos/app/command"

type Commands struct {
	RegistrarPago       command.RegistrarPago
	RegistrarPagoADeuda command.RegistarPagoADeuda
}
