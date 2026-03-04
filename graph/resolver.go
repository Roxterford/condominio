//go:generate go run ../internal/tools/autofilter/autofilter.go
//go:generate go run github.com/99designs/gqlgen generate
package graph

import (
	administracionService "github.com/Sanaruca/condominio/internal/administracion/service"
	usuarioService "github.com/Sanaruca/condominio/internal/usuarios/service"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require
// here.

type Resolver struct {
	Usuarios       *usuarioService.UsuarioService
	Administracion *administracionService.AdministracionService
}

func NewResolver(
	usuarioService *usuarioService.UsuarioService,
	administracionService *administracionService.AdministracionService,
) *Resolver {
	if usuarioService == nil {
		panic("usuarioService is required")
	}
	if administracionService == nil {
		panic("administracionService is required")
	}
	return &Resolver{
		Usuarios:       usuarioService,
		Administracion: administracionService,
	}
}
