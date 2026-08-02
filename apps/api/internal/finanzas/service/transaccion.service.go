package service

import (
	"github.com/Sanaruca/condominio/internal/core/common/quantity"
	"github.com/Sanaruca/condominio/internal/finanzas/app"
	"github.com/Sanaruca/condominio/internal/finanzas/app/command"
	"github.com/Sanaruca/condominio/internal/finanzas/app/query"
	"github.com/Sanaruca/condominio/internal/finanzas/models/transaccion"
	"github.com/Sanaruca/condominio/internal/unidades/models/unidad"
)

type TransaccionService struct {
	Commands app.Commands
	Queries  app.Queries
}

func New(
	repo transaccion.TransaccionRepository,
	factory *transaccion.TransaccionFactory,
	unidadRepo unidad.UnidadRepository,
	quantityFactory *quantity.QuantityFactory,
) *TransaccionService {
	return &TransaccionService{
		Commands: app.Commands{
			RegistrarTransaccion: command.NewRegistrarTransaccion(
				repo, factory, unidadRepo, quantityFactory,
			),
		},
		Queries: app.Queries{
			ObtenerMovimientos: query.NewObtenerMovimientos(repo),
		},
	}
}
