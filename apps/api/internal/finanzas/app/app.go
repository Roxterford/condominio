package app

import (
	"github.com/Sanaruca/condominio/internal/finanzas/app/command"
	"github.com/Sanaruca/condominio/internal/finanzas/app/query"
)

type Commands struct {
	RegistrarTransaccion command.RegistrarTransaccion
}

type Queries struct {
	ObtenerMovimientos query.ObtenerMovimientos
}
