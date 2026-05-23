package app

import (
	"github.com/Sanaruca/condominio/internal/pagos/app/command"
	"github.com/Sanaruca/condominio/internal/pagos/app/query"
)

type Commands struct {
	RegistrarPago             command.RegistrarPago
	ReprocesarPagosHuerfanos  command.ReprocesarPagosHuerfanos
}

type Queries struct {
	ObtenerPagos query.ObtenerPagos
}
