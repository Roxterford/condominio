package villa

import (
	"context"

	"github.com/Sanaruca/condominio/internal/core"
	"github.com/Sanaruca/condominio/internal/villas/models/villa/estadovilla"
)

type VillaRepository interface {
	Exists(ctx context.Context, villa int) (bool, core.Error)
	ObtenerTodas(ctx context.Context) ([]int, core.Error)
	ObtenerEstado(ctx context.Context, villa int) (estadovilla.EstadoDeVilla, core.Error)
}
