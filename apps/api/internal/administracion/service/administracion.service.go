package service

import (
	"github.com/Sanaruca/condominio/internal/administracion/app"
	"github.com/Sanaruca/condominio/internal/administracion/app/command"
	"github.com/Sanaruca/condominio/internal/administracion/app/query"
	"github.com/Sanaruca/condominio/internal/administracion/models/cuota"
	"github.com/Sanaruca/condominio/internal/administracion/models/deuda"
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
	proveedorRepository proveedor.ProveedorRepository,
	gastoRepository gasto.GastoRepository,
	cuotaRepository cuota.CuotaRepository,
	deudaRepository deuda.DeudaRepository,
	recaudacionFinder cuota.RecaudacionFinder,
	tasaService tasa.TasaService,
	proveedorFactory *proveedor.ProveedorFactory,
	emailFactory *common.EmailFactory,
	phoneFactory *common.PhoneFactory,
) *AdministracionService {

	registrarGasto := command.NewRegistrarGasto(gastoRepository, tasaService)
	registrarProveedor := command.NewRegistrarProveedor(
		proveedorRepository,
		proveedorFactory,
		emailFactory,
		phoneFactory,
	)

	return &AdministracionService{
		Queries: app.Queries{
			ObtenerProveedor:   query.NewObtenerProveedor(proveedorRepository),
			ObtenerProveedores: query.NewObtenerProveedores(proveedorRepository),
			ObtenerGastos:      query.NewObtenerGastos(gastoRepository),
			ObtenerCuotas:      query.NewObtenerCuotas(cuotaRepository),
			ObtenerCuota:       query.NewObtenerCuota(cuotaRepository),
			ObtenerRecaudacion: query.NewObtenerRecaudacion(recaudacionFinder),
		},
		Commands: app.Commands{
			RegistrarProveedor: registrarProveedor,
			RegistrarGasto:     registrarGasto,
			RegistrarGastoYProveedor: command.NewRegistrarGastoYProveedor(
				registrarProveedor,
				registrarGasto,
			),
			EliminarProveedor: command.NewEliminarProveedor(proveedorRepository),
		}}
}
