package graph

import (
	pagoService "github.com/Sanaruca/condominio/internal/pagos/service"
	usuarioService "github.com/Sanaruca/condominio/internal/usuarios/service"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require
// here.

type Resolver struct {
	Pagos    *pagoService.PagoService
	Usuarios *usuarioService.UsuarioService
}

func NewResolver(
	pagoService *pagoService.PagoService,
	usuarioService *usuarioService.UsuarioService,
) *Resolver {
	if pagoService == nil {
		panic("pagoService is required")
	}
	if usuarioService == nil {
		panic("usuarioService is required")
	}
	return &Resolver{
		Pagos:    pagoService,
		Usuarios: usuarioService,
	}
}
