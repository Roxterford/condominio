package query

import (
	"fmt"
	"time"

	"context"

	"github.com/Sanaruca/condominio/internal/core"
	"github.com/Sanaruca/condominio/internal/core/usecase"
	"github.com/Sanaruca/condominio/internal/services/tasa"
)

type ObtenerTasa usecase.WithOutput[int]

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

func (uc obtenerTasa) Exec(ctx context.Context, _ any) (int, core.Error) {

	now := time.Now()
	formatted := now.Format("2 January 2006 15:04:05")
	fmt.Println("time.Now(): ", formatted)

	_tasa, err := uc.tasaService.ObtenerTasaParaFecha(tasa.CambioOficial, time.Now())
	if err != nil {
		return 0, core.WrapError(err)
	}

	fmt.Println(_tasa)

	return _tasa.Valor, nil
}
