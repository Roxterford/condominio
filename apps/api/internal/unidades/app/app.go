package app

import "github.com/Sanaruca/condominio/internal/unidades/app/query"

type Queries struct {
	ObtenerUnidades                   query.ObtenerUnidades
	ObtenerUnidad                     query.ObtenerUnidad
	ObtenerUnidadPorCodigo            query.ObtenerUnidadPorCodigo
	ObtenerDeudasDeUnaUnidad          query.ObtenerDeudasDeUnaUnidadPorID
	ObtenerDeudasDeUnaUnidadPorCodigo query.ObtenerDeudasDeUnaUnidadPorCodigo
	ObtenerSujeto                     query.ObtenerSujeto
	ObtenerEstadisticas               query.ObtenerUnidadesEstadisticas
}
