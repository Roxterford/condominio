package app

import (
	"github.com/Sanaruca/condominio/internal/finanzas/app/command"
	"github.com/Sanaruca/condominio/internal/finanzas/app/query"
)

type Commands struct {
	RegistrarOperacion command.RegistrarOperacion
}

type Queries struct {
	ObtenerMovimientos query.ObtenerMovimientos
	ObtenerGastos      query.ObtenerGastos
}
