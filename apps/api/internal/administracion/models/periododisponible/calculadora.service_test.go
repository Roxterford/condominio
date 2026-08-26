package periododisponible

import (
	"testing"
	"time"

	"github.com/Sanaruca/condominio/internal/core/common/mes"
	"github.com/Sanaruca/condominio/internal/core/common/periodo"
)

func mustPeriodo(t *testing.T, m mes.Mes, anio int) periodo.Periodo {
	t.Helper()
	p, err := periodo.Nuevo(m, anio)
	if err != nil {
		t.Fatalf("periodo inválido: %v", err)
	}
	return p
}

func regla(t *testing.T, n int) ReglaDeEmision {
	t.Helper()
	r, err := Nuevo(n)
	if err != nil {
		t.Fatalf("regla inválida: %v", err)
	}
	return r
}

func TestEsValidoParaEmision_Duplicado(t *testing.T) {
	calc := NuevaCalculadoraDePeriodos()
	actual := periodo.DesdeTime(time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC))
	emitidos := []periodo.Periodo{mustPeriodo(t, mes.Julio, 2026)}

	err := calc.EsValidoParaEmision(
		mustPeriodo(t, mes.Julio, 2026), actual, regla(t, 1), emitidos, false,
	)
	if err == nil {
		t.Fatal("se esperaba error de duplicado")
	}
}

func TestEsValidoParaEmision_ContiguoAdelante(t *testing.T) {
	calc := NuevaCalculadoraDePeriodos()
	actual := periodo.DesdeTime(time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC))
	emitidos := []periodo.Periodo{mustPeriodo(t, mes.Julio, 2026)}

	if err := calc.EsValidoParaEmision(
		mustPeriodo(t, mes.Agosto, 2026), actual, regla(t, 1), emitidos, false,
	); err != nil {
		t.Fatalf("agosto debe ser válido contiguo: %v", err)
	}

	if err := calc.EsValidoParaEmision(
		mustPeriodo(t, mes.Septiembre, 2026), actual, regla(t, 1), emitidos, false,
	); err == nil {
		t.Fatal("septiembre debe ser rechazado por laguna")
	}
}

func TestEsValidoParaEmision_FueraDeVentana(t *testing.T) {
	calc := NuevaCalculadoraDePeriodos()
	actual := periodo.DesdeTime(time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC))
	emitidos := []periodo.Periodo{mustPeriodo(t, mes.Julio, 2026)}

	if err := calc.EsValidoParaEmision(
		mustPeriodo(t, mes.Octubre, 2026), actual, regla(t, 1), emitidos, false,
	); err == nil {
		t.Fatal("octubre debe ser rechazado por fuera de ventana")
	}

	// Extraordinaria relaja la ventana.
	if err := calc.EsValidoParaEmision(
		mustPeriodo(t, mes.Octubre, 2026), actual, regla(t, 1), emitidos, true,
	); err != nil {
		t.Fatalf("extraordinaria debe permitir octubre: %v", err)
	}
}

func TestEsValidoParaEmision_Bootstrap(t *testing.T) {
	calc := NuevaCalculadoraDePeriodos()
	actual := periodo.DesdeTime(time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC))

	if err := calc.EsValidoParaEmision(
		mustPeriodo(t, mes.Julio, 2026), actual, regla(t, 1), nil, false,
	); err != nil {
		t.Fatalf("bootstrap debe permitir julio: %v", err)
	}
}

func TestPeriodosDisponibles_Bootstrap(t *testing.T) {
	calc := NuevaCalculadoraDePeriodos()
	actual := periodo.DesdeTime(time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC))

	disp := calc.PeriodosDisponibles(actual, regla(t, 1), nil)
	if len(disp) != 3 {
		t.Fatalf("se esperaban 3 períodos (jul, ago, sep), got %d", len(disp))
	}
	primero := disp[0].Periodo()
	if primero.Mes() != mes.Julio || primero.Anio() != 2026 {
		t.Fatalf("el primero debe ser julio 2026, got %s", primero)
	}
}

func TestPeriodosDisponibles_LagunaYProyeccion(t *testing.T) {
	calc := NuevaCalculadoraDePeriodos()
	actual := periodo.DesdeTime(time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC))
	emitidos := []periodo.Periodo{
		mustPeriodo(t, mes.Enero, 2026),
		mustPeriodo(t, mes.Marzo, 2026), // feb faltante (laguna)
	}

	disp := calc.PeriodosDisponibles(actual, regla(t, 1), emitidos)
	// feb (laguna) + abr, may, jun, jul, ago, sep (proyección 1 mes)
	if len(disp) != 7 {
		t.Fatalf("se esperaban 7 períodos, got %d", len(disp))
	}
}
