package command

import (
	"errors"
	"time"

	"github.com/Sanaruca/condominio/internal/administracion"
	"github.com/Sanaruca/condominio/internal/core"
	"github.com/Sanaruca/condominio/internal/core/context"
	"github.com/Sanaruca/condominio/internal/core/usecase"
	"github.com/Sanaruca/condominio/internal/villas"
	"github.com/lucsky/cuid"
	"gorm.io/gorm"
)

type AplicarCuotaDTO struct {
	CuotaID string
	// TODO: Seria posible que se pueda aplicar la cuota a ciertas villas
	// Villas  []int
}

type AplicarCuota usecase.WithContextInput[context.BaseContext, AplicarCuotaDTO]

type aplicarCuota struct{}

func NewAplicarCuota() AplicarCuota {
	return &aplicarCuota{}
}

func (a *aplicarCuota) Exec(ctx context.BaseContext, dto AplicarCuotaDTO) (any, core.Error) {

	if dto.CuotaID == "" {
		return nil, core.NewInvalidArgumentError("Cuota es requerida")
	}

	_, err := gorm.G[administracion.Cuota](
		ctx.DB,
	).Select("id").
		Where("id = ?", dto.CuotaID).
		First(ctx)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, administracion.ErrCuotaNoEncontrada
	}

	if err != nil {
		return nil, core.WrapError(err)
	}

	var nvillas []int
	if err := ctx.DB.Model(new(villas.Villa)).Pluck("numero", &nvillas).Error; err != nil {
		return nil, core.WrapError(err)
	}

	if len(nvillas) < 1 {
		return nil, nil
	}

	deudas := make([]*villas.IDeuda, len(nvillas))
	for i, villa := range nvillas {
		deudas[i] = &villas.IDeuda{
			ID:            cuid.New(),
			Cuota:         dto.CuotaID,
			Villa:         villa,
			Registro:      time.Now(),
			Actualizacion: time.Now(),
		}
	}

	if err := ctx.DB.Create(&deudas).Error; err != nil {
		return nil, core.WrapError(err)
	}

	return nil, nil
}
