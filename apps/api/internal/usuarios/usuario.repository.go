package usuarios

import (
	"context"

	"github.com/Sanaruca/condominio/internal/core"
)

type UsuarioRepository interface {
	GetByEmail(ctx context.Context, email string) (*Usuario, core.Error)
	Guardar(ctx context.Context, usuario *Usuario) core.Error
}
