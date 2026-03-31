package tasa

import (
	"context"
	"time"
)

// TasaRepository define el contrato para acceder a datos de tasas de cambio.
// Es una interfaz de infraestructura/persistencia que representa diferentes fuentes
// de datos (base de datos local, APIs externas, servicios web, etc.).
//
// Responsabilidades:
// - Acceso directo a datos crudos de tasas
// - Conexión con fuentes externas (BCV, bancos, etc.)
// - Implementación de diferentes adaptadores (local, remoto, híbrido)
// - No contiene lógica de negocio ni caché
type TasaRepository interface {
	Nombre() string
	ObtenerTasaActual(tipo TipoDeCambio) (Tasa, error)
	ObtenerTasaPorFecha(tipo TipoDeCambio, fecha time.Time) (Tasa, error)
	ObtenerHistorico(tipo TipoDeCambio, desde, hasta time.Time) ([]Tasa, error)
}

// TasaCacheRepository define el contrato para el repositorio de caché de tasas
type TasaCacheRepository interface {
	Obtener(ctx context.Context, tipo TipoDeCambio, fecha time.Time) (Tasa, error)
	Guardar(ctx context.Context, tasa Tasa) error
}
