package gorm

import (
	"context"

	"gorm.io/gorm"

	"github.com/Sanaruca/condominio/internal/core"
	"github.com/Sanaruca/condominio/internal/core/common"
	"github.com/Sanaruca/condominio/internal/core/common/filter"
	"github.com/Sanaruca/condominio/internal/finanzas/models/operacion"
)

type GORMGastoFinder struct {
	db *gorm.DB
}

func NewGORMGastoFinder(db *gorm.DB) operacion.GastoFinder {
	return &GORMGastoFinder{db}
}

// Buscar implements [operacion.GastoFinder].
func (g *GORMGastoFinder) Buscar(
	ctx context.Context,
	filter filter.Clause,
	paginator common.Paginator,
) (*common.Paginated[operacion.Gasto], core.Error) {
	panic("unimplemented")
}
