package app

import (
	"github.com/Sanaruca/condominio/internal/administracion/app/command"
	"github.com/Sanaruca/condominio/internal/administracion/app/query"
)

type Queries struct {
	ObtenerProveedor   query.ObtenerProveedor
	ObtenerProveedores query.ObtenerProveedores
}

type Commands struct {
	RegistrarProveedor command.RegistrarProveedor
	EliminarProveedor  command.EliminarProveedor
	RegistrarGasto     command.RegistrarGasto
}
