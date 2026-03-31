//go:generate go run ../internal/tools/autofilter/autofilter.go
//go:generate go run github.com/99designs/gqlgen generate
package graph

import (
	administracionService "github.com/Sanaruca/condominio/internal/administracion/service"
	pagoService "github.com/Sanaruca/condominio/internal/pagos/service"
	sistemaService "github.com/Sanaruca/condominio/internal/sistema/service"
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
	Sistema        *sistemaService.SistemaService
}

func NewResolver(
	usuarioService *usuarioService.UsuarioService,
	administracionService *administracionService.AdministracionService,
	pagoService *pagoService.PagoService,
	sistemaService *sistemaService.SistemaService,
) *Resolver {
	if usuarioService == nil {
		panic("usuarioService is required")
	}
	if administracionService == nil {
		panic("administracionService is required")
	}
	if pagoService == nil {
		panic("pagoService is required")
	}
	if sistemaService == nil {
		panic("sistemaService is required")
	}
	return &Resolver{
		Usuarios:       usuarioService,
		Administracion: administracionService,
		Pagos:          pagoService,
		Sistema:        sistemaService,
	}
}
