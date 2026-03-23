package app

import (
	"github.com/Sanaruca/condominio/internal/administracion/app/command"
	"github.com/Sanaruca/condominio/internal/administracion/app/query"
	"github.com/Sanaruca/condominio/internal/administracion/models/proveedor"
)

type Queries struct {
	ObtenerProveedor   query.ObtenerProveedor
	ObtenerProveedores query.ObtenerProveedores
}

type Commands struct {
	RegistrarProveedor command.RegistrarProveedor
	EliminarProveedor  command.EliminarProveedor
}

func New(repo proveedor.ProveedorRepository) (Queries, Commands) {
	return Queries{
			ObtenerProveedor:   query.NewObtenerProveedor(repo),
			ObtenerProveedores: query.NewObtenerProveedores(repo),
		}, Commands{
			RegistrarProveedor: command.NewRegistrarProveedor(repo),
			EliminarProveedor:  command.NewEliminarProveedor(repo),
		}
}
