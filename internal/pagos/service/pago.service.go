package service

import (
	gormAdapter "github.com/Sanaruca/condominio/internal/pagos/adapters/gorm"
	redisAdapter "github.com/Sanaruca/condominio/internal/pagos/adapters/redis"
	"github.com/Sanaruca/condominio/internal/pagos/app"
	"github.com/Sanaruca/condominio/internal/pagos/app/command"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type PagoService struct {
	Commands app.Commands
}

func New(db *gorm.DB, redis_client *redis.Client) *PagoService {

	event_bus := redisAdapter.NewRedisEventBus(redis_client, "pagos")
	pago_repository := gormAdapter.NewPagoGORMRepository(db)

	registrarPago := command.NewRegistrarPago(pago_repository, nil, event_bus)

	return &PagoService{
		Commands: app.Commands{
			RegistrarPago: registrarPago,
		},
	}
}
