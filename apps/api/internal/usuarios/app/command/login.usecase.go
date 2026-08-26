package command

import (
	"fmt"
	"time"

	"github.com/Sanaruca/condominio/internal/core"
	"github.com/Sanaruca/condominio/internal/core/context"
	"github.com/Sanaruca/condominio/internal/core/envirotment"
	"github.com/Sanaruca/condominio/internal/core/exception"
	"github.com/Sanaruca/condominio/internal/core/usecase"
	"github.com/Sanaruca/condominio/internal/usuarios"
	"github.com/golang-jwt/jwt/v5"
)

type LoginDTO struct {
	Email    string
	Password string
}

type LoginCredentialsDTO struct {
	Token string `json:"token"`
}

type Login usecase.Handler[context.BaseContext, LoginDTO, *LoginCredentialsDTO]

func NewLogin(repo usuarios.UsuarioRepository) Login {
	return login{repo: repo}
}

type login struct {
	repo usuarios.UsuarioRepository // Inyección de dependencia
}

func (uc login) Exec(ctx context.BaseContext, input LoginDTO) (*LoginCredentialsDTO, core.Error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	usuario, err := uc.repo.GetByEmail(ctx, input.Email)

	error_credenciales := exception.New(exception.NOT_FOUND, "Credenciales inválidas")

	if err != nil {
		return nil, core.WrapError(err)
	}

	if usuario == nil {
		return nil, error_credenciales
	}

	if !usuario.ValidarPassword(input.Password) {
		return nil, error_credenciales
	}

	tokenString, err := uc.generateToken(usuario)

	if err != nil {
		fmt.Println(err, err != nil, err == nil, nil)
		return nil, err
	}

	return &LoginCredentialsDTO{Token: tokenString}, nil
}

func (dto LoginDTO) Validate() core.Error {

	if dto.Email == "" {
		return core.NewInvalidArgumentError("campo 'email' es requerido")
	}

	if dto.Password == "" {
		return core.NewInvalidArgumentError("campo 'password' es requerido")
	}

	return nil
}

// Helper privado para no ensuciar Exec
func (uc login) generateToken(u *usuarios.Usuario) (string, core.Error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss":   "condominio",
		"iat":   jwt.NewNumericDate(time.Now()),
		"ueid":  u.ID(),    // Usamos la Entidad
		"email": u.Email(), // Usamos la Entidad
		"exp":   jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
	})

	signedToken, err := token.SignedString([]byte(envirotment.GetSecretKey()))

	return signedToken, core.WrapError(err)
}
