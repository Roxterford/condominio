package app

import "github.com/Sanaruca/condominio/internal/administracion/app/query"

type Queries struct {
	ObtenerGastosSegunCuota query.ObtenerGastosSegunCuota
	ObtenerCuotas           query.ObtenerCuotas
}

type Commands struct {
}
