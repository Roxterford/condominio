package errors

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

func Wrap(err error) CoreError {
	if err != nil {

		if coreErr, ok := err.(CoreError); ok && coreErr.Code() != UNKNOWN {
			return coreErr
		}

		return &coreError{
			code:    UNKNOWN,
			message: err.Error(),
			cause:   err,
		}
	}

	return nil
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
