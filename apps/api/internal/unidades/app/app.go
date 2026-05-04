package app

import "github.com/Sanaruca/condominio/internal/unidades/app/query"

type Queries struct {
	ObtenerUnidades     query.ObtenerUnidades
	ObtenerSujeto       query.ObtenerSujeto
	ObtenerEstadisticas query.ObtenerUnidadesEstadisticas
}
