package ozzo

import (
	"fmt"
	"reflect"

	"github.com/Sanaruca/condominio/internal/core"
	"github.com/Sanaruca/condominio/internal/core/exception"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func init() {
	validation.ErrRequired = validation.ErrRequired.SetMessage("valor requerido")
	validation.ErrMinGreaterEqualThanRequired = validation.ErrMinGreaterEqualThanRequired.SetMessage(
		"valor debe ser mayor o igual a {{.threshold}}",
	)
	validation.ErrMinGreaterThanRequired = validation.ErrMinGreaterThanRequired.SetMessage(
		"valor debe ser mayor a {{.threshold}}",
	)
}

func OzzoErrorAdapter(dto core.Validable, err error) []core.Error {
	return ozzoErrorAdapter(dto, err, false)
}

func FirstOzzoErrorAdapter(dto core.Validable, err error) core.Error {
	errs := ozzoErrorAdapter(dto, err, true)
	if len(errs) == 0 {
		return nil
	}
	return errs[0]
}

func ozzoErrorAdapter(dto core.Validable, err error, firstOnly bool) []core.Error {
	if err == nil {
		return nil
	}

	var coreErrors []core.Error

	if validation_error, ok := err.(validation.Errors); ok {
		structName := reflect.TypeOf(dto).Name()
		for field, ferr := range validation_error {
			coreErrors = append(
				coreErrors,
				exception.NewInvalidArgumentError(
					"%s", fmt.Sprintf("`%s.%s`: %s", structName, field, ferr.Error()),
				),
			)
			if firstOnly {
				break
			}
		}
	} else {
		// Error genérico
		coreErrors = append(coreErrors, exception.WithCause(exception.NewValidationError("%s", err.Error()), err))
	}

	return coreErrors
}
