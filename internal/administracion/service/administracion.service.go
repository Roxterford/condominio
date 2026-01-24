package service

import (
	"github.com/Sanaruca/condominio/internal/administracion/app"
	"github.com/Sanaruca/condominio/internal/administracion/app/query"
)

type AdministracionService struct {
	Queries  app.Queries
	Commands app.Commands
}

func New() *AdministracionService {

	obtenerGastosSegunCuota := query.NewObtenerGastosSegunCuota()

	return &AdministracionService{
		Queries: app.Queries{
			ObtenerGastosSegunCuota: obtenerGastosSegunCuota,
		},
		Commands: app.Commands{},
	}
}
