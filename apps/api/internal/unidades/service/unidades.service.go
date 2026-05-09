package service

import (
	"github.com/Sanaruca/condominio/internal/administracion/models/deuda"
	"github.com/Sanaruca/condominio/internal/unidades/app"
	"github.com/Sanaruca/condominio/internal/unidades/app/query"
	"github.com/Sanaruca/condominio/internal/unidades/models/sujeto"
	"github.com/Sanaruca/condominio/internal/unidades/models/unidad"
)

type UnidadesService struct {
	Queries app.Queries
}

func NewUnidadesService(
	unidadRepository unidad.UnidadRepository,
	sujetoRepository sujeto.SujetoRepository,
	deudaRepository deuda.DeudaRepository,
	unidadEstadisticasFinder unidad.EstadisticasFinder,
) *UnidadesService {
	return &UnidadesService{
		Queries: app.Queries{
			ObtenerUnidades:        query.NewObtenerUnidades(unidadRepository),
			ObtenerUnidad:          query.NewObtenerUnidad(unidadRepository),
			ObtenerUnidadPorCodigo: query.NewObtenerUnidadPorCodigo(unidadRepository),
			ObtenerSujeto:          query.NewObtenerSujeto(sujetoRepository),
			ObtenerEstadisticas: query.NewObtenerUnidadesEstadisticas(
				unidadEstadisticasFinder,
			),
			ObtenerDeudasDeUnaUnidad: query.NewObtenerDeudasDeUnaUnidadPorID(
				deudaRepository,
			),
			ObtenerDeudasDeUnaUnidadPorCodigo: query.NewObtenerDeudasDeUnaUnidadPorCodigo(
				deudaRepository,
			),
		},
	}
}
