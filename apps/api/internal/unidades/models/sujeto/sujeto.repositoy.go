package sujeto

import (
	"context"

	"github.com/Sanaruca/condominio/internal/core"
)

type SujetoRepository interface {
	ObtenerPorID(ctx context.Context, id SujetoID) (Sujeto, core.Error)
	Guardar(ctx context.Context, sujeto Sujeto) core.Error
	ExistsDocumento(ctx context.Context, documento string) (bool, core.Error)
	ExistsEmail(ctx context.Context, email string) (bool, core.Error)
}
