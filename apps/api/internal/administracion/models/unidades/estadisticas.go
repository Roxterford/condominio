package unidades

import (
	"context"

	"github.com/Sanaruca/condominio/internal/core"
	"github.com/Sanaruca/condominio/internal/core/common/quantity"
)

type UnidadesEstadisticas struct {
	TotalUnidades         int
	UnidadesActivas       int
	UnidadesInhabitadas   int
	UnidadesExentas       int
	UnidadesEnLitigio     int
	UnidadesSuspendidas   int
	UnidadesPreventa      int
	UnidadesConPendientes int
	UnidadesSolventes     int
	TotalPendiente        quantity.Quantity
	TotalAsignado         quantity.Quantity
}

type UnidadesEstadisticasFinder interface {
	Obtener(ctx context.Context) (*UnidadesEstadisticas, core.Error)
}
