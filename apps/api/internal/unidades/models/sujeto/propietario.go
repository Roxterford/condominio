package sujeto

type Propietario interface {
	Sujeto
	Contacto() Persona
}

type propietario struct {
	Persona
}

func (p propietario) Contacto() Persona {
	return p.Persona
}
