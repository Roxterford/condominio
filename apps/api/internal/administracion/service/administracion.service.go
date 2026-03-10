package service

import (
	"github.com/Sanaruca/condominio/internal/administracion/app"
)

type AdministracionService struct {
	Queries  app.Queries
	Commands app.Commands
}

func New() *AdministracionService {

	return &AdministracionService{
		Queries:  app.Queries{},
		Commands: app.Commands{},
	}
}
