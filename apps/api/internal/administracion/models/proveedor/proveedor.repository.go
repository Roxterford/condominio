package proveedor

import (
	"context"

	"github.com/Sanaruca/condominio/internal/core"
	"github.com/Sanaruca/condominio/internal/core/common/filter"
)

type ProveedorRepository interface {
	Guardar(ctx context.Context, proveedor *Proveedor) core.Error
	ObtenerPorID(ctx context.Context, id string) (*Proveedor, core.Error)
	// Deprecated: Por razones de rendimiento
	ObtenerTodos(ctx context.Context) ([]*Proveedor, core.Error)
	Obtener(ctx context.Context, filter filter.Clause) ([]Proveedor, core.Error)
	Actualizar(ctx context.Context, proveedor *Proveedor) core.Error
	Eliminar(ctx context.Context, id string) core.Error
	ExistePorRif(ctx context.Context, rif string) (bool, core.Error)
	ExistePorID(ctx context.Context, id string) (bool, core.Error)
}
