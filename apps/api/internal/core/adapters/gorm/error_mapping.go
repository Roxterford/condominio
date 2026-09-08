package gormAdapter

import (
	"gorm.io/gorm"

	"github.com/Sanaruca/condominio/internal/core/exception"
)

func init() {
	exception.RegisterMapping(gorm.ErrRecordNotFound, exception.Mapping{
		Code:    exception.NOT_FOUND,
		Message: "Recurso no encontrado",
	})

	exception.RegisterMapping(gorm.ErrDuplicatedKey, exception.Mapping{
		Code:    exception.CONFLICT,
		Message: "El recurso ya existe",
	})
}
