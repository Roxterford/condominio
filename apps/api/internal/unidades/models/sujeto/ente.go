package sujeto

// Ente representa a un ente, como una empresa, sociedad o cooperativa,
// que puede ser un sujeto de interés del condominio.
type Ente struct {
	sujeto
	razon_social  string
	representante Persona
}

// RazonSocial devuelve la razón social del ente.
func (e Ente) RazonSocial() string { return e.razon_social }

// Representante devuelve la persona que representa al ente.
func (e Ente) Representante() Persona { return e.representante }
