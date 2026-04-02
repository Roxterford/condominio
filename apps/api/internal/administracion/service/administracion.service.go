package service

import (
	"github.com/Sanaruca/condominio/internal/administracion/app"
	"github.com/Sanaruca/condominio/internal/administracion/app/command"
	"github.com/Sanaruca/condominio/internal/administracion/app/query"
	"github.com/Sanaruca/condominio/internal/administracion/models/gasto"
	"github.com/Sanaruca/condominio/internal/administracion/models/proveedor"
	"github.com/Sanaruca/condominio/internal/core/common"
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
	proveedorFactory *proveedor.ProveedorFactory,
	emailFactory *common.EmailFactory,
	phoneFactory *common.PhoneFactory,
) *AdministracionService {

	registrarGasto := command.NewRegistrarGasto(gastoRepository, tasaService)
	registrarProveedor := command.NewRegistrarProveedor(
		proveedorRepositoy,
		proveedorFactory,
		emailFactory,
		phoneFactory,
	)

	return &AdministracionService{
		app.Queries{
			ObtenerProveedor:   query.NewObtenerProveedor(proveedorRepositoy),
			ObtenerProveedores: query.NewObtenerProveedores(proveedorRepositoy),
			ObtenerGastos:      query.NewObtenerGastos(gastoRepository),
		},
		app.Commands{
			RegistrarProveedor: registrarProveedor,
			RegistrarGasto:     registrarGasto,
			RegistrarGastoYProveedor: command.NewRegistrarGastoYProveedor(
				registrarProveedor,
				registrarGasto,
			),
			EliminarProveedor: command.NewEliminarProveedor(proveedorRepositoy),
		}}
}
