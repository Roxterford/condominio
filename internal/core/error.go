package core

import (
	"fmt"

	"github.com/Sanaruca/condominio/internal/core/errors"
)

type Error = errors.CoreError

func NewError(code errors.CoreErrorCode, message string) Error {
	return errors.New(code, message)
}

func WrapError(err error) Error {
	return errors.Wrap(err)
}

func NewValidationError(message string) Error {
	return errors.NewValidationError(message)
}

func NewInvalidArgumentError(message string, msgArgs ...any) Error {
	return errors.NewInvalidArgumentError(fmt.Sprintf(message, msgArgs...))
}
