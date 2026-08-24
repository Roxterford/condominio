package core

import (
	"fmt"

	"github.com/Sanaruca/condominio/internal/core/exception"
)

type Error = exception.CoreError

var ErrInternal = exception.New(exception.INTERNAL, "Ocurrió un error inesperado")

// TODO: add msgArgs to all error functions

func NewError(code exception.CoreErrorCode, message string) Error {
	return exception.New(code, "%s", message)
}

func WrapError(err error) Error {
	return exception.Wrap(err)
}

func NewValidationError(message string) Error {
	return exception.NewValidationError("%s", message)
}

func NewInvalidArgumentError(message string, msgArgs ...any) Error {
	return exception.NewInvalidArgumentError("%s", fmt.Sprintf(message, msgArgs...))
}
