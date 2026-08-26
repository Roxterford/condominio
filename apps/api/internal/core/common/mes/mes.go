package mes

import "github.com/Sanaruca/condominio/internal/core"

type Mes int

const (
	Enero Mes = iota + 1
	Febrero
	Marzo
	Abril
	Mayo
	Junio
	Julio
	Agosto
	Septiembre
	Octubre
	Noviembre
	Diciembre
)

var nombres = map[Mes]string{
	Enero:      "Enero",
	Febrero:    "Febrero",
	Marzo:      "Marzo",
	Abril:      "Abril",
	Mayo:       "Mayo",
	Junio:      "Junio",
	Julio:      "Julio",
	Agosto:     "Agosto",
	Septiembre: "Septiembre",
	Octubre:    "Octubre",
	Noviembre:  "Noviembre",
	Diciembre:  "Diciembre",
}

func (m Mes) Validate() core.Error {
	if m < 1 || m > 12 {
		return core.NewValidationError("mes invalido")
	}
	return nil
}
func (m Mes) Value() int {
	return int(m)
}

func (m Mes) String() string {
	return nombres[m]
}
