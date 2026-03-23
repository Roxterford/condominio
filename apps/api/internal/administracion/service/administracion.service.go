package service

import (
	"github.com/Sanaruca/condominio/internal/administracion/app"
	"github.com/Sanaruca/condominio/internal/administracion/models/proveedor"
)

type AdministracionService struct {
	Queries  app.Queries
	Commands app.Commands
}

func New(proveedorRepo proveedor.ProveedorRepository) *AdministracionService {
	queries, commands := app.New(proveedorRepo)

	return &AdministracionService{
		Queries:  queries,
		Commands: commands,
	}
}
