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

func (m Mes) Validate() core.Error {
	if m < 1 || m > 12 {
		return core.NewValidationError("mes invalido")
	}
	return nil
}

func (m Mes) Value() int {
	return int(m)
}
