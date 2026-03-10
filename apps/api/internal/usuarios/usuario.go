package usuarios

import (
	"github.com/Sanaruca/condominio/internal/core"
	"github.com/Sanaruca/condominio/internal/core/errors"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUsuarioNoEncontrado = errors.New(errors.NOT_FOUND, "usuario no encontrado")
)

// TODO: Plantear nuevos tipos de usuario (admin, propietario, etc.)

type Usuario struct {
	id       string
	email    string
	password Password
	nombre   string
	apellido string
}

func (u *Usuario) ID() string         { return u.id }
func (u *Usuario) Email() string      { return u.email }
func (u *Usuario) Password() Password { return u.password }
func (u *Usuario) Nombre() string     { return u.nombre }
func (u *Usuario) Apellido() string   { return u.apellido }

// TODO: Aqui podriamos añadir un hasher como dependencia
type UsuarioFactory struct{}

func NewFactory() UsuarioFactory {
	return UsuarioFactory{}
}

func (f UsuarioFactory) Assemble(
	id, email, password, nombre, apellido string,
) (*Usuario, core.Error) {
	// Validaciones de integridad (No son reglas de negocio, son de estructura)
	if id == "" {
		return nil, core.NewInvalidArgumentError("identidad faltante en base de datos")
	}
	if email == "" {
		return nil, core.NewInvalidArgumentError("email faltante en base de datos")
	}

	// Rehidratamos el objeto
	return &Usuario{
		id:       id,
		email:    email,
		password: Password{hash: password},
		nombre:   nombre,
		apellido: apellido,
	}, nil
}

func (u *Usuario) ValidarPassword(password string) bool {
	return u.password.Check(password)
}

// Password es nuestro Value Object
type Password struct {
	hash string
}

// NewPassword crea un hash a partir de un string plano (regla de dominio)
func NewPassword(plainText string) (Password, core.Error) {
	if len(plainText) < 8 {
		return Password{}, errors.New(
			errors.INVALID_ARGUMENT,
			"la contraseña debe tener al menos 8 caracteres",
		)
	}

	bytes, err := bcrypt.GenerateFromPassword([]byte(plainText), bcrypt.DefaultCost)
	if err != nil {
		return Password{}, core.WrapError(err)
	}

	return Password{hash: string(bytes)}, nil
}

// Check valida si la contraseña coincide
func (p Password) Check(plainText string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(p.hash), []byte(plainText))
	return err == nil
}

// String devuelve el hash (útil para la capa de infraestructura/DB)
func (p Password) String() string {
	return p.hash
}
