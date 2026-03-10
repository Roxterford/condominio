package core

import (
	"fmt"

	"github.com/Sanaruca/condominio/internal/core/errors"
)

type Error = errors.CoreError

// TODO: add msgArgs to all error functions

func NewError(code errors.CoreErrorCode, message string) Error {
	return errors.New(code, "%s", message)
}

func WrapError(err error) Error {
	return errors.Wrap(err)
}

func NewValidationError(message string) Error {
	return errors.NewValidationError("%s", message)
}

func NewInvalidArgumentError(message string, msgArgs ...any) Error {
	return errors.NewInvalidArgumentError("%s", fmt.Sprintf(message, msgArgs...))
}
