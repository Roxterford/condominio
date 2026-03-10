package utils

import (
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func GetFirstOzzoValidationError(err error) error {
	if validationErrors, ok := err.(validation.Errors); ok {
		for key, err := range validationErrors {
			return fmt.Errorf(err.Error(), key)
		}
	}

	return err
}
