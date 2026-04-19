package service

import (
	"github.com/Sanaruca/condominio/internal/unidades/app"
	"github.com/Sanaruca/condominio/internal/unidades/app/query"
	"github.com/Sanaruca/condominio/internal/unidades/models/unidad"
)

type UnidadesService struct {
	Queries app.Queries
}

func NewUnidadesService(unidadRepository unidad.UnidadRepository) *UnidadesService {
	return &UnidadesService{
		Queries: app.Queries{
			ObtenerUnidades: query.NewObtenerUnidades(unidadRepository),
		},
	}
}
