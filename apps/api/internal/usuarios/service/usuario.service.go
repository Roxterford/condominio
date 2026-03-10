package service

import (
	"github.com/Sanaruca/condominio/internal/usuarios"
	"github.com/Sanaruca/condominio/internal/usuarios/app"
	"github.com/Sanaruca/condominio/internal/usuarios/app/command"
)

type UsuarioService struct {
	Commands app.Commands
}

func New(repo usuarios.UsuarioRepository) *UsuarioService {

	login := command.NewLogin(repo)
	return &UsuarioService{
		Commands: app.Commands{
			Login: login,
		},
	}
}
