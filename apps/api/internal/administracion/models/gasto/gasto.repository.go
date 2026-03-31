package gasto

import (
	"context"

	"github.com/Sanaruca/condominio/internal/core"
)

type GastoRepository interface {
	Guardar(ctx context.Context, gasto Gasto) (*GastoID, core.Error)
	ObtenerPorID(ctx context.Context, id GastoID) (*Gasto, core.Error)
}
