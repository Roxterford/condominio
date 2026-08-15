package service

import (
	"github.com/Sanaruca/condominio/internal/administracion/app"
	"github.com/Sanaruca/condominio/internal/administracion/app/command"
	"github.com/Sanaruca/condominio/internal/administracion/app/query"
	"github.com/Sanaruca/condominio/internal/administracion/models/cuota"
	"github.com/Sanaruca/condominio/internal/administracion/models/proveedor"
	"github.com/Sanaruca/condominio/internal/core/common"
)

type AdministracionService struct {
	Queries  app.Queries
	Commands app.Commands
}

func New(
	proveedorRepository proveedor.ProveedorRepository,
	cuotaRepository cuota.CuotaRepository,
	recaudacionFinder cuota.RecaudacionFinder,
	proveedorFactory *proveedor.ProveedorFactory,
	emailFactory *common.EmailFactory,
	phoneFactory *common.PhoneFactory,
	cuotaFactory *cuota.CuotaFactory,
) *AdministracionService {

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
			ObtenerCuotas:      query.NewObtenerCuotas(cuotaRepository),
			ObtenerCuota:       query.NewObtenerCuota(cuotaRepository),
			ObtenerRecaudacion: query.NewObtenerRecaudacion(recaudacionFinder),
		},
		Commands: app.Commands{
			RegistrarProveedor: registrarProveedor,
			EliminarProveedor:  command.NewEliminarProveedor(proveedorRepository),
		}}
}
