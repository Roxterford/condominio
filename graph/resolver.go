package graph

import (
	administracionService "github.com/Sanaruca/condominio/internal/administracion/service"
	pagoService "github.com/Sanaruca/condominio/internal/pagos/service"
	usuarioService "github.com/Sanaruca/condominio/internal/usuarios/service"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require
// here.

type Resolver struct {
	Pagos          *pagoService.PagoService
	Usuarios       *usuarioService.UsuarioService
	Administracion *administracionService.AdministracionService
}

func NewResolver(
	pagoService *pagoService.PagoService,
	usuarioService *usuarioService.UsuarioService,
	administracionService *administracionService.AdministracionService,
) *Resolver {
	if pagoService == nil {
		panic("pagoService is required")
	}
	if usuarioService == nil {
		panic("usuarioService is required")
	}
	if administracionService == nil {
		panic("administracionService is required")
	}
	return &Resolver{
		Pagos:          pagoService,
		Usuarios:       usuarioService,
		Administracion: administracionService,
	}
}
