package service

import (
	"github.com/Sanaruca/condominio/internal/core/common/quantity"
	"github.com/Sanaruca/condominio/internal/finanzas/app"
	"github.com/Sanaruca/condominio/internal/finanzas/app/command"
	"github.com/Sanaruca/condominio/internal/finanzas/app/query"
	"github.com/Sanaruca/condominio/internal/finanzas/models/operacion"
	"github.com/Sanaruca/condominio/internal/unidades/models/unidad"
)

type FinanzaService struct {
	Commands app.Commands
	Queries  app.Queries
}

func New(
	repo operacion.OperacionRepository,
	gastoFinder operacion.GastoFinder,
	factory *operacion.OperacionFactory,
	unidadRepo unidad.UnidadRepository,
	quantityFactory *quantity.QuantityFactory,
) *FinanzaService {
	return &FinanzaService{
		Commands: app.Commands{
			RegistrarOperacion: command.NewRegistrarTransaccion(
				repo, factory, unidadRepo, quantityFactory,
			),
		},
		Queries: app.Queries{
			ObtenerMovimientos: query.NewObtenerMovimientos(repo),
			ObtenerGastos:      query.NewObtenerGastos(gastoFinder),
		},
	}
}
