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

// Add suma dos cantidades si pertenecen a la misma escala.
func (q Quantity) Add(other Quantity) (Quantity, error) {
	if q.scale != other.scale {
		return Quantity{}, ErrScaleMismatch
	}
	return Quantity{value: q.value + other.value, scale: q.scale}, nil
}

// HappyAdd suma dos cantidades ignorando la escala del otro operando.
func (q Quantity) HappyAdd(other Quantity) Quantity {
	return Quantity{value: q.value + other.value, scale: q.scale}
}

// Sub resta dos cantidades si pertenecen a la misma escala.
func (q Quantity) Sub(other Quantity) (Quantity, error) {
	if q.scale != other.scale {
		return Quantity{}, ErrScaleMismatch
	}
	return Quantity{value: q.value - other.value, scale: q.scale}, nil
}

// HappySub resta dos cantidades ignorando la escala del otro operando.
func (q Quantity) HappySub(other Quantity) Quantity {
	return Quantity{value: q.value - other.value, scale: q.scale}
}

// Mul multiplica dos cantidades si pertenecen a la misma escala.
// El resultado se reescala para mantener la misma escala original.
func (q Quantity) Mul(other Quantity) (Quantity, error) {
	if q.scale != other.scale {
		return Quantity{}, ErrScaleMismatch
	}
	divisor := int64(math.Pow10(q.scale))
	result := int64(math.Round(float64(q.value*other.value) / float64(divisor)))
	return Quantity{value: result, scale: q.scale}, nil
}

// HappyMul multiplica dos cantidades ignorando la escala del otro operando.
func (q Quantity) HappyMul(other Quantity) Quantity {
	divisor := int64(math.Pow10(q.scale))
	result := int64(math.Round(float64(q.value*other.value) / float64(divisor)))
	return Quantity{value: result, scale: q.scale}
}

// Div divide dos cantidades si pertenecen a la misma escala.
// Retorna error si el divisor es cero.
func (q Quantity) Div(other Quantity) (Quantity, error) {
	if q.scale != other.scale {
		return Quantity{}, ErrScaleMismatch
	}
	if other.value == 0 {
		return Quantity{}, errors.NewInvalidArgumentError("cannot divide by zero")
	}
	multiplier := int64(math.Pow10(q.scale))
	result := int64(math.Round(float64(q.value*multiplier) / float64(other.value)))
	return Quantity{value: result, scale: q.scale}, nil
}

// HappyDiv divide dos cantidades ignorando la escala del otro operando.
// Retorna un Quantity vacío si el divisor es cero.
func (q Quantity) HappyDiv(other Quantity) Quantity {
	if other.value == 0 {
		return Quantity{scale: q.scale}
	}
	multiplier := int64(math.Pow10(q.scale))
	result := int64(math.Round(float64(q.value*multiplier) / float64(other.value)))
	return Quantity{value: result, scale: q.scale}
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
