package service

import (
	"github.com/Sanaruca/condominio/internal/services/tasa"
	"github.com/Sanaruca/condominio/internal/sistema/app"
	"github.com/Sanaruca/condominio/internal/sistema/app/query"
)

type SistemaService struct {
	Queries app.Queries
}

func New(tasa_service tasa.TasaService) *SistemaService {
	return &SistemaService{
		Queries: app.Queries{
			ObtenerTasa: query.NewObtenerTasa(tasa_service),
		},
	}
}
