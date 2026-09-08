package exception

import (
	"errors"
	"fmt"
)

type CoreErrorCode string

const (
	UNKNOWN          CoreErrorCode = "UNKNOWN"
	VALIDATION       CoreErrorCode = "VALIDATION"
	INVALID_ARGUMENT CoreErrorCode = "INVALID_ARGUMENT"
	NOT_FOUND        CoreErrorCode = "NOT_FOUND"
	CONFLICT         CoreErrorCode = "CONFLICT"
	INTERNAL         CoreErrorCode = "INTERNAL"
	UNAUTHORIZED     CoreErrorCode = "UNAUTHORIZED"
	FORBIDDEN        CoreErrorCode = "FORBIDEN"
)

type CoreError interface {
	error
	Code() CoreErrorCode
	Message() string
	Cause() error
}

func New(code CoreErrorCode, message string, msgArgs ...interface{}) CoreError {
	return &coreError{
		code:    code,
		message: fmt.Sprintf(message, msgArgs...),
	}
}

func WithCause(err CoreError, cause error) CoreError {

	if err == nil {
		return nil
	}

	return &coreError{
		code:    err.Code(),
		message: err.Error(),
		cause:   cause,
	}
}

func NewValidationError(message string, msgArgs ...interface{}) CoreError {
	return New(VALIDATION, message, msgArgs...)
}

func NewInvalidArgumentError(message string, msgArgs ...interface{}) CoreError {
	return New(INVALID_ARGUMENT, message, msgArgs...)
}

// Mapping define cómo transformar un error externo (SQL, redis, etc.) en un
// error de dominio que el cliente puede ver. El detalle interno nunca se
// expone: se mantiene como causa y solo el mensaje mapeado llega al cliente.
type Mapping struct {
	Code CoreErrorCode
	// Message es una plantilla estilo fmt.Sprintf. Si no tiene verbos de
	// formato, se usa tal cual.
	Message string
	// MessageArgs es opcional. Recibe el error original y devuelve los
	// argumentos para renderizar Message.
	MessageArgs func(err error) []interface{}
}

// errorMappings es el registro de mapeos de errores conocidos. Solo se escribe
// en init() de cada adapter y se lee en runtime, así que no requiere locks.
var errorMappings = map[error]Mapping{}

// RegisterMapping registra cómo mapear un error sentinel conocido a un error
// de dominio. Debe llamarse desde init() de los adapters, nunca en runtime.
func RegisterMapping(target error, mapping Mapping) {
	errorMappings[target] = mapping
}

// TODO: en el futuro no queremos exponer ningun detalle de errores internos
// o desconocidos
func Wrap(err error) CoreError {
	if err == nil {
		return nil
	}

	if coreErr, ok := err.(CoreError); ok && coreErr.Code() != UNKNOWN {
		return coreErr
	}

	for target, mapping := range errorMappings {
		if errors.Is(err, target) {
			var message string
			if mapping.MessageArgs != nil {
				message = fmt.Sprintf(mapping.Message, mapping.MessageArgs(err)...)
			} else {
				message = mapping.Message
			}

			return &coreError{
				code:    mapping.Code,
				message: message,
				cause:   err,
			}
		}
	}

	return &coreError{
		code:    INTERNAL,
		message: "Ocurrió un error inesperado",
		cause:   err,
	}
}

type coreError struct {
	code    CoreErrorCode
	message string
	cause   error
}

func (e *coreError) Error() string {
	return fmt.Sprintf("[%s]: %s", e.code, e.message)
}

func (e *coreError) Code() CoreErrorCode {
	return e.code
}

func (e *coreError) Message() string {
	return e.message
}

func (e *coreError) Cause() error {
	return e.cause
}

func (e *coreError) WithCause(cause error) CoreError {
	e.cause = cause
	return e
}

func Is(err, target error) bool {
	return errors.Is(err, target)
}
