package sujeto

// Ente representa a un ente, como una empresa, sociedad o cooperativa,
// que puede ser un sujeto de interés del condominio.
type Ente struct {
	sujeto
	razon_social  string
	representante Persona
}

// AsEnte implements [Sujeto].
func (e *Ente) AsEnte() *Ente {
	return e
}

// AsPersona implements [Sujeto].
func (e *Ente) AsPersona() *Persona {
	return nil
}

// AsPropietario implements [Sujeto].
func (e *Ente) AsPropietario() Propietario {
	return propietario{
		Persona: e.representante,
	}
}

func (e Ente) RazonSocial() string    { return e.razon_social }
func (e Ente) Representante() Persona { return e.representante }
func (e Ente) Contacto() Persona      { return e.representante }
