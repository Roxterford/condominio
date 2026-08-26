package sujeto

type Titular interface {
	Sujeto
	Contacto() Persona
}

type titular struct {
	Sujeto
	contacto Persona
}

func (p titular) Contacto() Persona {
	return p.contacto
}
