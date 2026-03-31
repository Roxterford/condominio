package tasa

import (
	"time"

	"github.com/Sanaruca/condominio/internal/core/errors"
)

type TipoDeCambio string

const (
	CambioOficial  TipoDeCambio = "oficial"
	CambioParalelo TipoDeCambio = "paralelo"
	CambioPromedio TipoDeCambio = "promedio"
)

type Tasa struct {
	Valor  int
	Fuente string
	Fecha  time.Time
	Tipo   TipoDeCambio
	// TODO: las implementaciones usan VES/USD, VES/EUR, etc. Posible normalización futura
	Moneda string
}

var (
	ErrTasaNoEncontrada = errors.New(errors.NOT_FOUND, "Tasa de cambio no encontrada")
	ErrTasaNoDisponible = errors.New(
		errors.INTERNAL,
		"Tasa de cambio no disponible para la fecha solicitada",
	)
	ErrSinConexion = errors.New(
		errors.INTERNAL,
		"Sin conexión al servicio de tasas de cambio",
	)
	ErrAdaptadorInvalido = errors.New(errors.INVALID_ARGUMENT, "Adaptador de cambio no válido")
)

// TasaService define el contrato para el servicio de dominio de tasas de cambio.
// Es una interfaz de aplicación que contiene lógica de negocio y reglas de dominio.
//
// Responsabilidades:
// - Orquestar múltiples repositorios
// - Implementar caché para optimizar rendimiento
// - Aplicar reglas de negocio (ej: usar tasa oficial para pagos)
// - Proporcionar una API simplificada al resto de la aplicación
//
// Diferencia clave con TasaRepository:
// - TasaRepository = CÓMO obtener datos (infraestructura)
// - TasaService = QUÉ hacer con los datos (dominio)
type TasaService interface {
	// ObtenerTasaParaFecha obtiene la tasa de cambio correspondiente al tipo de cambio y fecha solicitada.
	//
	// Parámetros:
	// - tipo: el tipo de cambio solicitado (oficial, paralelo, promedio)
	// - fecha: la fecha deseada para obtener la tasa de cambio
	//
	// Retorna:
	// - tasa: la tasa de cambio obtenida
	// - error: si ocurrió algún error durante la búsqueda de la tasa de cambio
	ObtenerTasaParaFecha(tipo TipoDeCambio, fecha time.Time) (Tasa, error)

	// ObtenerTasaParaPago obtiene la tasa de cambio correspondiente al tipo de cambio oficial y fecha del pago.
	//
	// Parámetros:
	// - fechaPago: la fecha del pago para obtener la tasa de cambio
	//
	// Retorna:
	// - tasa: la tasa de cambio obtenida
	// - error: si ocurrió algún error durante la búsqueda de la tasa de cambio
	ObtenerTasaParaPago(fechaPago time.Time) (Tasa, error)
}
