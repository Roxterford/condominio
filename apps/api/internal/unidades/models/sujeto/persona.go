package sujeto

type Persona struct {
	sujeto
	nombres   string
	apellidos string
}

func (p Persona) Nombres() string   { return p.nombres }
func (p Persona) Apellidos() string { return p.apellidos }

func (p Persona) AsPropietario() Propietario { return &propietario{p} }
func (p Persona) AsPersona() *Persona        { return &p }
func (p Persona) AsEnte() *Ente              { return nil }
