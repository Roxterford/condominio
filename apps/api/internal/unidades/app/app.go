package app

import "github.com/Sanaruca/condominio/internal/unidades/app/query"

type Queries struct {
	ObtenerUnidades        query.ObtenerUnidades
	ObtenerUnidad          query.ObtenerUnidad
	ObtenerUnidadPorCodigo query.ObtenerUnidadPorCodigo
	ObtenerSujeto          query.ObtenerSujeto
	ObtenerEstadisticas    query.ObtenerUnidadesEstadisticas
}
