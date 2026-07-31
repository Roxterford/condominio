package app

import (
	"github.com/Sanaruca/condominio/internal/transacciones/app/command"
	"github.com/Sanaruca/condominio/internal/transacciones/app/query"
)

type Commands struct {
	RegistrarTransaccion command.RegistrarTransaccion
}

type Queries struct {
	ObtenerMovimientos query.ObtenerMovimientos
}
