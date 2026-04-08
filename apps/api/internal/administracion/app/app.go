package app

import (
	"github.com/Sanaruca/condominio/internal/administracion/app/command"
	"github.com/Sanaruca/condominio/internal/administracion/app/query"
)

type Queries struct {
	ObtenerProveedor   query.ObtenerProveedor
	ObtenerProveedores query.ObtenerProveedores
	ObtenerGastos      query.ObtenerGastos
	ObtenerCuotas      query.ObtenerCuotas
}

type Commands struct {
	RegistrarProveedor       command.RegistrarProveedor
	RegistrarGasto           command.RegistrarGasto
	RegistrarGastoYProveedor command.RegistrarGastoYProveedor
	EliminarProveedor        command.EliminarProveedor
}
