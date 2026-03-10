package gorm

import "github.com/Sanaruca/condominio/internal/usuarios"

type Usuario struct {
	ID       string
	Email    string
	Password string
	Nombre   string
	Apellido string
}

// mapToUsuarioTable convierte la cabecera del usuario del dominio al modelo de GORM
func mapToUsuarioTable(u *usuarios.Usuario) *Usuario {
	return &Usuario{
		ID:       u.ID(),
		Email:    u.Email(),
		Password: u.Password().String(),
		Nombre:   u.Nombre(),
		Apellido: u.Apellido(),
	}
}
