package periodo

import (
	"strconv"
	"time"

	"github.com/Sanaruca/condominio/internal/core"
	"github.com/Sanaruca/condominio/internal/core/common/mes"
)

type Periodo struct {
	mes  mes.Mes
	anio int
}

func Nuevo(m mes.Mes, anio int) (Periodo, core.Error) {
	if err := m.Validate(); err != nil {
		return Periodo{}, err
	}
	if anio < 2000 {
		return Periodo{}, core.NewValidationError("el año debe ser mayor a 2000")
	}
	return Periodo{mes: m, anio: anio}, nil
}

func Assemble(m mes.Mes, anio int) Periodo {
	return Periodo{mes: m, anio: anio}
}

func DesdeTime(t time.Time) Periodo {
	return Periodo{mes: mes.Mes(t.Month()), anio: t.Year()}
}

func (p Periodo) Mes() mes.Mes { return p.mes }
func (p Periodo) Anio() int    { return p.anio }

func (p Periodo) String() string {
	return p.mes.String() + " " + strconv.Itoa(p.anio)
}

func (p Periodo) totalMeses() int {
	return p.anio*12 + int(p.mes) - 1
}

func desdeTotalMeses(t int) Periodo {
	anio := t / 12
	mesN := t%12 + 1
	return Periodo{mes: mes.Mes(mesN), anio: anio}
}

func (p Periodo) Siguiente() Periodo {
	return desdeTotalMeses(p.totalMeses() + 1)
}

func (p Periodo) Anterior() Periodo {
	return desdeTotalMeses(p.totalMeses() - 1)
}

func (p Periodo) Desplazar(n int) Periodo {
	return desdeTotalMeses(p.totalMeses() + n)
}

func (p Periodo) EsIgual(o Periodo) bool {
	return p.totalMeses() == o.totalMeses()
}

func (p Periodo) EsPosteriorA(o Periodo) bool {
	return p.totalMeses() > o.totalMeses()
}

func (p Periodo) EsAnteriorA(o Periodo) bool {
	return p.totalMeses() < o.totalMeses()
}

func (p Periodo) DistanciaEnMesesA(o Periodo) int {
	d := p.totalMeses() - o.totalMeses()
	if d < 0 {
		d = -d
	}
	return d
}
