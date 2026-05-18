package query

import (
	"context"
	"time"

	"github.com/Sanaruca/condominio/internal/core"
	"github.com/Sanaruca/condominio/internal/core/common/mes"
	"github.com/Sanaruca/condominio/internal/core/common/quantity"
	"github.com/Sanaruca/condominio/internal/core/usecase"
	"github.com/Sanaruca/condominio/internal/services/tasa"
)

type ObtenerTasaDTO struct {
	Dia  int
	Mes  mes.Mes
	Anio int

	fecha time.Time
}

type ObtenerTasaResult struct {
	Valor quantity.Quantity
	Fecha time.Time
}

type ObtenerTasa usecase.Handler[context.Context, *ObtenerTasaDTO, *ObtenerTasaResult]

type obtenerTasa struct {
	tasaService tasa.TasaService
}

func NewObtenerTasa(tasaService tasa.TasaService) ObtenerTasa {
	if tasaService == nil {
		panic("tasaService is required")
	}
	return obtenerTasa{
		tasaService: tasaService,
	}
}

func (uc obtenerTasa) Exec(
	ctx context.Context,
	input *ObtenerTasaDTO,
) (*ObtenerTasaResult, core.Error) {

	input.Validate()

	_tasa, err := uc.tasaService.ObtenerTasaParaFecha(tasa.CambioOficial, input.fecha)
	if err != nil {
		return nil, core.WrapError(err)
	}

	return &ObtenerTasaResult{
		Valor: _tasa.Valor,
		Fecha: _tasa.Fecha,
	}, nil
}

func (input *ObtenerTasaDTO) Validate() core.Error {
	now := time.Now()

	// 1. Si el input es nil, no puedes simplemente reasignarlo así (no afectaría al puntero original fuera)
	// Pero asumiendo que el caller lo maneja, inicializamos con los valores de 'now'
	if input.Anio <= 0 {
		input.Anio = now.Year()
	}

	// 2. Validar Mes (usando tu tipo mes.Mes)
	// Asumo que mes.Mes tiene un método IsValid() o lo comparas con constantes
	if err := input.Mes.Validate(); err != nil {
		input.Mes = mes.Mes(int(now.Month()))
	}

	// 3. Validar Día
	// No todos los meses tienen 31 días. Usamos la lógica de time.Date para validar
	// Si el día es menor a 1 o mayor a lo que permite ese mes/año, default al día actual
	if input.Dia < 1 || input.Dia > 31 {
		input.Dia = now.Day()
	}

	// 4. Construir la fecha final y verificar si es una fecha real (ej. 31 de febrero)
	// Al usar time.Date, Go normaliza fechas (ej. 31 de Abril se convierte en 1 de Mayo)
	// Para cumplir tu regla de "si está mal -> actual", verificamos la validez:
	t := time.Date(input.Anio, time.Month(input.Mes), input.Dia, 0, 0, 0, 0, time.Local)

	// Si Go tuvo que "normalizar" la fecha (porque el usuario metió un día inexistente)
	// o si la fecha es en el futuro (si es que no permites tasas futuras),
	// podrías forzar el fallback aquí.
	if t.Month() != time.Month(input.Mes) {
		input.fecha = now
		input.Anio, input.Mes, input.Dia = now.Year(), mes.Mes(int(now.Month())), now.Day()
	} else {
		input.fecha = t
	}

	return nil
}
