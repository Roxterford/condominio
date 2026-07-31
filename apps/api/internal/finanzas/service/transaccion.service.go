package service

import (
	"github.com/Sanaruca/condominio/internal/core/common/quantity"
	"github.com/Sanaruca/condominio/internal/services/tasa"
	"github.com/Sanaruca/condominio/internal/transacciones/app"
	"github.com/Sanaruca/condominio/internal/transacciones/app/command"
	"github.com/Sanaruca/condominio/internal/transacciones/app/query"
	"github.com/Sanaruca/condominio/internal/transacciones/models/transaccion"
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
	tasaService tasa.TasaService,
	quantityFactory *quantity.QuantityFactory,
) *TransaccionService {
	return &TransaccionService{
		Commands: app.Commands{
			RegistrarTransaccion: command.NewRegistrarTransaccion(
				repo, factory, unidadRepo, tasaService, quantityFactory,
			),
		},
		Queries: app.Queries{
			ObtenerMovimientos: query.NewObtenerMovimientos(repo),
		},
	}
}
