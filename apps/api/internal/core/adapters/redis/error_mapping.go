package redis

import (
	"github.com/redis/go-redis/v9"

	"github.com/Sanaruca/condominio/internal/core/exception"
)

func init() {
	exception.RegisterMapping(redis.Nil, exception.Mapping{
		Code:    exception.NOT_FOUND,
		Message: "Recurso no encontrado",
	})
}
