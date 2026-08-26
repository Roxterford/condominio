package sujeto

import (
	"context"

	"github.com/Sanaruca/condominio/internal/core"
)

type SujetoRepository interface {
	ObtenerPorID(ctx context.Context, id SujetoID) (Sujeto, core.Error)
}
