package command

import (
	"time"

	"github.com/Sanaruca/condominio/internal/core"
	"github.com/Sanaruca/condominio/internal/core/context"
	"github.com/Sanaruca/condominio/internal/core/envirotment"
	"github.com/Sanaruca/condominio/internal/core/errors"
	"github.com/Sanaruca/condominio/internal/core/usecase"
	"github.com/Sanaruca/condominio/internal/usuarios"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type LoginDTO struct {
	Email    string
	Password string
}

type LoginCredentialsDTO struct {
	Token string `json:"token"`
}

type Login usecase.Handler[context.BaseContext, LoginDTO, *LoginCredentialsDTO]

func NewLogin() Login {
	return login{}
}

type login struct {
}

func (uc login) Exec(ctx context.BaseContext, input LoginDTO) (*LoginCredentialsDTO, core.Error) {

	if err := input.Validate(); err != nil {
		return nil, err
	}

	usuario, err := gorm.G[usuarios.Usuario](
		ctx.DB,
	).Select("id", "email", "password").
		Where("email = ?", input.Email).
		First(ctx)

	error_credenciales := errors.New(errors.NOT_FOUND, "Credenciales invalidas")

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, error_credenciales
	}

	if err != nil {
		return nil, core.WrapError(err)
	}

	if usuario.Password != input.Password {
		return nil, error_credenciales
	}

	// Se estan especificando como claims los definidos en https://www.iana.org/assignments/jwt/jwt.xhtml
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss":   "condominio",
		"iat":   jwt.NewNumericDate(time.Now()),
		"ueid":  usuario.ID,
		"email": usuario.Email,
		// TODO: considere usar tokens de corta duracion y un refresh token
		"exp": jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // 24h
	})

	tokenString, err := token.SignedString([]byte(envirotment.GetSecretKey()))
	if err != nil {
		return nil, core.WrapError(err)
	}

	return &LoginCredentialsDTO{
		Token: tokenString,
	}, nil
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
