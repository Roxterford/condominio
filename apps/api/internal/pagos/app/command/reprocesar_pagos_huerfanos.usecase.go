package command

import (
	"github.com/Sanaruca/condominio/internal/core"
	"github.com/Sanaruca/condominio/internal/core/context"
	"github.com/Sanaruca/condominio/internal/core/usecase"
	"github.com/Sanaruca/condominio/internal/pagos/models/pago"
)

type ReprocesarPagosHuerfanosResultado struct {
	Total      int      `json:"total"`
	Procesados int      `json:"procesados"`
	Errores    int      `json:"errores"`
	Detalles   []string `json:"detalles,omitempty"`
}

type ReprocesarPagosHuerfanos usecase.Handler[context.BaseContext, any, *ReprocesarPagosHuerfanosResultado]

type reprocesarPagosHuerfanos struct {
	pagos       pago.PagoRepository
	aplicarPago AplicarPago
}

func NewReprocesarPagosHuerfanos(
	pagoRepository pago.PagoRepository,
	aplicarPago AplicarPago,
) ReprocesarPagosHuerfanos {
	return &reprocesarPagosHuerfanos{
		pagos:       pagoRepository,
		aplicarPago: aplicarPago,
	}
}

func (uc *reprocesarPagosHuerfanos) Exec(
	ctx context.BaseContext,
	_ any,
) (*ReprocesarPagosHuerfanosResultado, core.Error) {
	pagos, err := uc.pagos.ObtenerPagosSinDestinos(ctx)

	if err != nil {
		return nil, err
	}

	resultado := &ReprocesarPagosHuerfanosResultado{
		Total:    len(pagos),
		Detalles: make([]string, 0),
	}

	for _, p := range pagos {
		_, execErr := uc.aplicarPago.Exec(ctx, AplicarPagoDTO{
			PagoID: p.ID(),
		})

		if execErr != nil {
			resultado.Errores++
			resultado.Detalles = append(
				resultado.Detalles,
				"pago "+p.ID()+": "+execErr.Error(),
			)
		} else {
			resultado.Procesados++
		}
	}

	return resultado, nil
}
