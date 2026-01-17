package service

import (
	"github.com/Sanaruca/condominio/internal/usuarios/app"
	"github.com/Sanaruca/condominio/internal/usuarios/app/command"
)

type UsuarioService struct {
	Commands app.Commands
}

func New() *UsuarioService {
	login := command.NewLogin()
	return &UsuarioService{
		Commands: app.Commands{
			Login: login,
		},
	}
}
