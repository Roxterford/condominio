package common

import (
	"math"

	"github.com/Sanaruca/condominio/internal/core/common/quantity"
)

// Percentage es un Value Object que representa una proporción.
// Generalmente usa una escala mayor (ej. 4) para permitir
// precisiones como 5.25% (0.0525).
type Percentage struct {
	q quantity.Quantity
}

func (p Percentage) Value() int64 {
	return p.q.Value()
}

// NewPercentage crea un porcentaje con la escala definida.
func NewPercentage(value int64, scale int) Percentage {
	return Percentage{quantity.New(value, scale)}
}

// DisplayString devuelve la representación visual con el símbolo %.
// Si la magnitud es 525 y escala 4, representa 5.25%.
func (p Percentage) DisplayString() string {
	// Para mostrarlo como "5.25" en lugar de "0.0525",
	// multiplicamos mentalmente por 100 (reducimos la escala en 2)
	displayScale := p.q.Scale() - 2
	if displayScale < 0 {
		displayScale = 0
	}

	// Reutilizamos la lógica de Ratio pero ajustada a formato humano de %
	tempQ := quantity.New(p.q.Value(), displayScale)
	return tempQ.Ratio() + "%"
}

// ApplyTo calcula el porcentaje sobre una cantidad dada.
// Fórmula: (Quantity * Percentage.value) / 10^Percentage.scale
func (p Percentage) ApplyTo(target quantity.Quantity) quantity.Quantity {
	// Calculamos la nueva magnitud
	// Nota: Esto puede requerir manejar desbordamientos en casos extremos
	newvalue := (target.Value() * p.q.Value()) / int64(math.Pow10(p.q.Scale()))

	return quantity.New(newvalue, target.Scale())
}
