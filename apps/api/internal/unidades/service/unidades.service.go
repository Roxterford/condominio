package service

import (
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
	unidadEstadisticasFinder unidad.EstadisticasFinder,
) *UnidadesService {
	return &UnidadesService{
		Queries: app.Queries{
			ObtenerUnidades:     query.NewObtenerUnidades(unidadRepository),
			ObtenerSujeto:       query.NewObtenerSujeto(sujetoRepository),
			ObtenerEstadisticas: query.NewObtenerUnidadesEstadisticas(unidadEstadisticasFinder),
		},
	}
}
