package hybrid

import (
	"time"

	"github.com/Sanaruca/condominio/internal/services/tasa"
)

// HybridTasaRepository es un adaptador híbrido que combina múltiples fuentes de tasas de cambio.
// Implementa una estrategia de fallback donde primero intenta obtener datos de una fuente local
// y si falla, recurre secuencialmente a fuentes remotas hasta encontrar una respuesta exitosa.
//
// Estrategia de funcionamiento:
// 1. Prioridad siempre a la fuente local (base de datos local)
// 2. Si la fuente local falla, intenta con cada fuente remota en orden
// 3. Retorna la primera respuesta exitosa encontrada
// 4. Si todas las fuentes fallan, retorna ErrTasaNoDisponible
//
// Las tasas obtenidas de la fuente local se marcan con el sufijo " (Local)"
// para distinguirlas de las tasas provenientes de fuentes remotas.
type HybridTasaRepository struct {
	// local es el repositorio principal de tasas (usualmente base de datos local)
	// Tiene prioridad sobre todas las demás fuentes
	local tasa.TasaRepository

	// remotos es una lista de repositorios externos (APIs, servicios web, etc.)
	// Se utilizan como fallback cuando la fuente local no está disponible
	remotos []tasa.TasaRepository
}

func NewHybridTasaRepository(
	local tasa.TasaRepository,
	remotos []tasa.TasaRepository,
) *HybridTasaRepository {
	return &HybridTasaRepository{
		local:   local,
		remotos: remotos,
	}
}

func (a *HybridTasaRepository) Nombre() string {
	return "hybrid"
}

func (a *HybridTasaRepository) ObtenerTasaActual(tipo tasa.TipoDeCambio) (tasa.Tasa, error) {
	_tasa, err := a.local.ObtenerTasaActual(tipo)
	if err == nil {
		_tasa.Fuente = _tasa.Fuente + " (Local)"
		return _tasa, nil
	}

	for _, remoto := range a.remotos {
		_tasa, err := remoto.ObtenerTasaActual(tipo)
		if err == nil {
			return _tasa, nil
		}
	}

	return tasa.Tasa{}, tasa.ErrTasaNoDisponible
}

func (a *HybridTasaRepository) ObtenerTasaPorFecha(
	tipo tasa.TipoDeCambio,
	fecha time.Time,
) (tasa.Tasa, error) {
	_tasa, err := a.local.ObtenerTasaPorFecha(tipo, fecha)
	if err == nil {
		_tasa.Fuente = _tasa.Fuente + " (Local)"
		return _tasa, nil
	}

	for _, remoto := range a.remotos {
		_tasa, err := remoto.ObtenerTasaPorFecha(tipo, fecha)
		if err == nil {
			return _tasa, nil
		}
	}

	return tasa.Tasa{}, tasa.ErrTasaNoDisponible
}

func (a *HybridTasaRepository) ObtenerHistorico(
	tipo tasa.TipoDeCambio,
	desde, hasta time.Time,
) ([]tasa.Tasa, error) {
	localTasas, err := a.local.ObtenerHistorico(tipo, desde, hasta)
	if err == nil && len(localTasas) > 0 {
		for i := range localTasas {
			localTasas[i].Fuente = localTasas[i].Fuente + " (Local)"
		}
		return localTasas, nil
	}

	for _, remoto := range a.remotos {
		tasas, err := remoto.ObtenerHistorico(tipo, desde, hasta)
		if err == nil && len(tasas) > 0 {
			return tasas, nil
		}
	}

	return nil, tasa.ErrTasaNoDisponible
}
