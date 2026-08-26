package tasa

import (
	"context"
	"time"

	"github.com/Sanaruca/condominio/internal/core"
)

type tasaService struct {
	repo      TasaRepository
	cacheRepo TasaCacheRepository
}

func NewTasaService(repo TasaRepository, cacheRepo TasaCacheRepository) TasaService {
	if repo == nil {
		panic("repo cannot be nil")
	}
	return &tasaService{
		repo:      repo,
		cacheRepo: cacheRepo,
	}
}

func (s *tasaService) ObtenerTasaParaFecha(tipo TipoDeCambio, fecha time.Time) (*Tasa, core.Error) {
	ctx := context.Background()

	if s.cacheRepo != nil {
		cached, err := s.cacheRepo.Obtener(ctx, tipo, fecha)
		if err == nil && cached.Valor.Value() > 0 {
			return &cached, nil
		}
	}

	tasa, err := s.repo.ObtenerTasaPorFecha(tipo, fecha)
	if err != nil {
		return nil, core.WrapError(err)
	}

	if s.cacheRepo != nil && tasa.Valor.Value() > 0 {
		_ = s.cacheRepo.Guardar(ctx, tasa)
	}

	return &tasa, nil
}

func (s *tasaService) ObtenerTasaParaPago(fechaPago time.Time) (*Tasa, core.Error) {
	return s.ObtenerTasaParaFecha(CambioOficial, fechaPago)
}

func (s *tasaService) ObtenerTasaActual(tipo TipoDeCambio) (*Tasa, core.Error) {
	return s.ObtenerTasaParaFecha(tipo, time.Now())
}
