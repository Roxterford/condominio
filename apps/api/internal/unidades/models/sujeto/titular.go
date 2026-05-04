package sujeto

type Titular interface {
	Sujeto
	Contacto() Persona
}

type titular struct {
	Persona
}

func (p titular) Contacto() Persona {
	return p.Persona
}
