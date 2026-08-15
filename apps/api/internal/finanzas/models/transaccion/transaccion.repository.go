package transaccion

import (
	"context"

	"github.com/Sanaruca/condominio/internal/core"
)

type TransaccionRepository interface {
	Guardar(ctx context.Context, tx *Transaccion) core.Error
	ObtenerPorID(ctx context.Context, id string) (*Transaccion, core.Error)
}
