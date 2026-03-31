package service

import (
	"github.com/Sanaruca/condominio/internal/administracion/app"
	"github.com/Sanaruca/condominio/internal/administracion/app/command"
	"github.com/Sanaruca/condominio/internal/administracion/app/query"
	"github.com/Sanaruca/condominio/internal/administracion/models/gasto"
	"github.com/Sanaruca/condominio/internal/administracion/models/proveedor"
	"github.com/Sanaruca/condominio/internal/services/tasa"
)

type AdministracionService struct {
	Queries  app.Queries
	Commands app.Commands
}

func New(
	proveedorRepositoy proveedor.ProveedorRepository,
	gastoRepository gasto.GastoRepository,
	tasaService tasa.TasaService,
) *AdministracionService {
	return &AdministracionService{
		app.Queries{
			ObtenerProveedor:   query.NewObtenerProveedor(proveedorRepositoy),
			ObtenerProveedores: query.NewObtenerProveedores(proveedorRepositoy),
		},
		app.Commands{
			RegistrarProveedor: command.NewRegistrarProveedor(proveedorRepositoy),
			EliminarProveedor:  command.NewEliminarProveedor(proveedorRepositoy),
			RegistrarGasto:     command.NewRegistrarGasto(gastoRepository, tasaService),
		}}
}
