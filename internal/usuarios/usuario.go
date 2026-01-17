package usuarios

import "github.com/Sanaruca/condominio/internal/core/errors"

var (
	ErrUsuarioNoEncontrado = errors.New(errors.NOT_FOUND, "usuario no encontrado")
)

type Usuario struct {
	ID       string
	Email    string
	Password string
	Nombre   string
	Apellido string
}

type UsuarioPayload struct {
	ID    string
	Email string
}
