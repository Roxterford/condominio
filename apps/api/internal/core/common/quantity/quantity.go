package quantity

import (
	"fmt"
	"math"

	"github.com/Sanaruca/condominio/internal/core/errors"
)

const DEFAULT_SCALE = 2

var ErrScaleMismatch = errors.NewInvalidArgumentError(
	"cannot operate on quantities with different scales",
)

// Quantity es un Value Object que representa una magnitud con escala fija.
// Es inmutable y agnóstico al contexto (dinero, %, masa, etc).
type Quantity struct {
	value int64 // El valor entero bruto
	scale int   // El exponente de base 10 (ej: 2 para centésimas)
}

func (q Quantity) Value() int64 {
	return q.value
}

func (q Quantity) Float() float64 {
	divisor := math.Pow10(q.scale)
	return float64(q.value) / divisor
}

func (q Quantity) Scale() int {
	return q.scale
}

// New crea una cantidad basada en su magnitud entera y su escala.
// Ejemplo: New(3992, 5) representa 0.03992
func New(value int64, scale int) Quantity {

	if scale < 0 {
		scale = DEFAULT_SCALE
	}

	return Quantity{
		value: value,
		scale: scale,
	}
}

// FromFloat convierte un float a Quantity con una precisión específica.
// Útil para entradas de usuario o APIs externas.
func FromFloat(f float64, scale int) Quantity {
	if scale < 0 {
		scale = DEFAULT_SCALE
	}
	multiplier := math.Pow10(scale)
	return Quantity{
		value: int64(math.Round(f * multiplier)),
		scale: scale,
	}
}

// Add suma dos magnitudes si pertenecen a la misma escala.
func (q Quantity) Add(other Quantity) (Quantity, error) {
	if q.scale != other.scale {
		return Quantity{}, ErrScaleMismatch
	}
	return Quantity{value: q.value + other.value, scale: q.scale}, nil
}

// String devuelve la representación decimal humana.
func (q Quantity) String() string {
	if q.scale <= 0 {
		return fmt.Sprintf("%d", q.value)
	}

	divisor := int64(math.Pow10(q.scale))
	beforeDot := q.value / divisor
	afterDot := q.value % divisor

	if afterDot < 0 {
		afterDot = -afterDot
	}

	format := fmt.Sprintf("%%d.%%0%dd", q.scale)
	return fmt.Sprintf(format, beforeDot, afterDot)
}
