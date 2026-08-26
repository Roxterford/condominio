package periododisponible

import (
	"github.com/Sanaruca/condominio/internal/core"
	"github.com/Sanaruca/condominio/internal/core/common/periodo"
)

var (
	ErrPeriodoDuplicado = core.NewValidationError(
		"ya existe una cuota para el período indicado",
	)
	ErrPeriodoFueraDeVentana = core.NewValidationError(
		"el período excede la ventana de emisión futura permitida",
	)
	ErrPeriodoNoContiguo = core.NewValidationError(
		"el período no es contiguo al historial de emisión (se detectaría una laguna)",
	)
)

type CalculadoraDePeriodos struct{}

func NuevaCalculadoraDePeriodos() *CalculadoraDePeriodos {
	return &CalculadoraDePeriodos{}
}

// PeriodosDisponibles determina los períodos en los que aún se puede emitir
// una cuota, considerando el historial ya emitido, el período actual y la
// regla de ventana futura. Incluye los meses faltantes en el historial
// (lagunas) y la proyección contigua hacia adelante.
func (c *CalculadoraDePeriodos) PeriodosDisponibles(
	actual periodo.Periodo,
	regla ReglaDeEmision,
	emitidos []periodo.Periodo,
) []PeriodoDisponible {

	if len(emitidos) == 0 {
		return rango(
			actual.Anterior(),
			actual.Desplazar(regla.MaximoMesesFuturo()),
			"Período disponible para emisión",
		)
	}

	ultimo := maxPeriodo(emitidos)
	primer := minPeriodo(emitidos)

	disponibles := make([]PeriodoDisponible, 0)

	// Lagunas: meses faltantes dentro del historial ya iniciado.
	for p := primer; !p.EsPosteriorA(ultimo); p = p.Siguiente() {
		if !contiene(emitidos, p) {
			disponibles = append(
				disponibles,
				NuevoPeriodoDisponible(p, "Período faltante en el historial (laguna)", false),
			)
		}
	}

	// Proyección hacia adelante, contigua al último emitido.
	limite := actual.Desplazar(regla.MaximoMesesFuturo())
	for p := ultimo.Siguiente(); !p.EsPosteriorA(limite); p = p.Siguiente() {
		disponibles = append(
			disponibles,
			NuevoPeriodoDisponible(p, "Período disponible para emisión", false),
		)
	}

	return disponibles
}

// EsValidoParaEmision valida que un período candidato puede registrarse como
// cuota, aplicando los invariantes de negocio:
//   - sin duplicidad respecto al historial emitido;
//   - dentro de la ventana futura (relajada para cuotas extraordinarias);
//   - contiguo al historial para evitar lagunas.
func (c *CalculadoraDePeriodos) EsValidoParaEmision(
	candidato periodo.Periodo,
	actual periodo.Periodo,
	regla ReglaDeEmision,
	emitidos []periodo.Periodo,
	esExtraordinario bool,
) core.Error {

	if contiene(emitidos, candidato) {
		return ErrPeriodoDuplicado
	}

	if esExtraordinario {
		// Las cuotas extraordinarias tienen respaldo (proyecto/justificación),
		// por lo que no se restringen por ventana futura ni contigüidad.
		return nil
	}

	limite := actual.Desplazar(regla.MaximoMesesFuturo())
	if candidato.EsPosteriorA(limite) {
		return ErrPeriodoFueraDeVentana
	}

	if len(emitidos) == 0 {
		return nil
	}

	ultimo := maxPeriodo(emitidos)
	primer := minPeriodo(emitidos)

	switch {
	case candidato.EsPosteriorA(ultimo):
		if candidato.DistanciaEnMesesA(ultimo) != 1 {
			return ErrPeriodoNoContiguo
		}
	case candidato.EsAnteriorA(primer):
		if candidato.DistanciaEnMesesA(primer) != 1 {
			return ErrPeriodoNoContiguo
		}
	default:
		return ErrPeriodoNoContiguo
	}

	return nil
}

func rango(desde, hasta periodo.Periodo, motivo string) []PeriodoDisponible {
	if desde.EsPosteriorA(hasta) {
		return []PeriodoDisponible{}
	}
	out := make([]PeriodoDisponible, 0)
	for p := desde; !p.EsPosteriorA(hasta); p = p.Siguiente() {
		out = append(out, NuevoPeriodoDisponible(p, motivo, false))
	}
	return out
}

func maxPeriodo(ps []periodo.Periodo) periodo.Periodo {
	m := ps[0]
	for _, p := range ps[1:] {
		if p.EsPosteriorA(m) {
			m = p
		}
	}
	return m
}

func minPeriodo(ps []periodo.Periodo) periodo.Periodo {
	m := ps[0]
	for _, p := range ps[1:] {
		if p.EsAnteriorA(m) {
			m = p
		}
	}
	return m
}

func contiene(ps []periodo.Periodo, p periodo.Periodo) bool {
	for _, x := range ps {
		if x.EsIgual(p) {
			return true
		}
	}
	return false
}
