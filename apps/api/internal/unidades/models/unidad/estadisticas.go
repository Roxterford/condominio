package unidad

import (
	"context"

	"github.com/Sanaruca/condominio/internal/core"
	"github.com/Sanaruca/condominio/internal/core/common/quantity"
)

type Estadisticas struct {
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

type EstadisticasFinder interface {
	Obtener(ctx context.Context) (*Estadisticas, core.Error)
}
