package quantity

// QuantityFactory centraliza la creación de cantidades bajo reglas específicas.
type QuantityFactory struct {
	defaultScale int
}

// NewFactory crea una nueva instancia de la factoría con una escala base.
func NewFactory(scale int) *QuantityFactory {
	return &QuantityFactory{defaultScale: scale}
}

// New crea una Quantity validando las reglas de negocio iniciales.
// Se usa para "Nuevas" instancias en el dominio (ej. desde un input de usuario).
func (f *QuantityFactory) New(value int64) Quantity {
	// Aquí podrías añadir validaciones de dominio, como no permitir negativos

	return New(value, f.defaultScale)
}

// Assemble reconstruye una Quantity usando la escala por defecto.
// Útil cuando tu DB es "muda" y confías en la definición actual del dominio.
func (f *QuantityFactory) Assemble(value int64) Quantity {
	return New(value, f.defaultScale)
}

// AssembleWithScale permite reconstruir una Quantity con una escala explícita.
// Ideal para migraciones, datos externos o cuando la DB provee precisión.
func (f *QuantityFactory) AssembleWithScale(value int64, scale int) Quantity {
	return New(value, scale)
}
