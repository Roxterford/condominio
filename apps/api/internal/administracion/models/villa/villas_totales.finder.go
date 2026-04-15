package villa

import (
	"context"

	"github.com/Sanaruca/condominio/internal/core"
	"github.com/Sanaruca/condominio/internal/core/common/quantity"
)

type VillasTotales struct {
	TotalVillas         int
	VillasActivas       int
	VillasInhabitadas   int
	VillasExentas       int
	VillasEnLitigio     int
	VillasSuspendidas   int
	VillasPreventa      int
	VillasConPendientes int
	VillasSolventes     int
	TotalPendiente      quantity.Quantity
	TotalAsignado       quantity.Quantity
}

type VillasTotalesFinder interface {
	ObtenerVillasTotales(ctx context.Context) (*VillasTotales, core.Error)
}
