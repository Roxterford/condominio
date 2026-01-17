package command

import (
	"time"

	"github.com/Sanaruca/condominio/internal/administracion"
	"github.com/Sanaruca/condominio/internal/administracion/types/tipodecuota"
	"github.com/Sanaruca/condominio/internal/core"
	"github.com/Sanaruca/condominio/internal/core/context"
	"github.com/Sanaruca/condominio/internal/core/usecase"
	"github.com/lucsky/cuid"
	"gorm.io/gorm"
)

type RegistrarCuotaDTO struct {
	Mes   int
	Anio  int
	Monto int
	Tipo  tipodecuota.TipoDeCuota
}

type RegistrarCuota usecase.WithContextInput[context.AdminContext, RegistrarCuotaDTO]

type registrarCuota struct {
}

func NewRegistrarCuota() RegistrarCuota {
	return &registrarCuota{}
}

func (r *registrarCuota) Exec(ctx context.AdminContext, dto RegistrarCuotaDTO) (any, core.Error) {

	if err := dto.Validate(); err != nil {
		return nil, err
	}

	gorm.G[administracion.Cuota](ctx.DB).Create(ctx, &administracion.Cuota{
		ID:            cuid.New(),
		Monto:         dto.Monto,
		Mes:           dto.Mes,
		Anio:          dto.Anio,
		Registro:      time.Now(),
		Actualizacion: time.Now(),
	})

	return nil, nil
}

func (dto RegistrarCuotaDTO) Validate() core.Error {

	if err := dto.Tipo.Validate(); err != nil {
		return err
	}

	if dto.Mes <= 0 || dto.Mes > 12 {
		return core.NewInvalidArgumentError("Mes '%d' no es valido", dto.Mes)
	}

	if dto.Monto <= 0 {
		return core.NewInvalidArgumentError("Monto '%d' no es valido", dto.Monto)
	}

	return nil
}
