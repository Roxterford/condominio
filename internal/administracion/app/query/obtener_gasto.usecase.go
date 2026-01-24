package query

import (
	"errors"
	"time"

	"github.com/Sanaruca/condominio/internal/administracion"
	"github.com/Sanaruca/condominio/internal/core"
	"github.com/Sanaruca/condominio/internal/core/context"
	"github.com/Sanaruca/condominio/internal/core/usecase"
	"gorm.io/gorm"
)

type ObtenerGastoDTO struct {
	GastoID string
}

type ObtenerGasto usecase.Handler[context.BaseContext, ObtenerGastoDTO, *administracion.Gasto]

func NewObtenerGasto() ObtenerGasto {
	return &obtenerGasto{}
}

type obtenerGasto struct{}

func (uc *obtenerGasto) Exec(
	ctx context.BaseContext,
	dto ObtenerGastoDTO,
) (*administracion.Gasto, core.Error) {

	gasto_map := map[string]any{}

	err := ctx.DB.Table((*administracion.Gasto).TableName(&administracion.Gasto{})).
		Where("id = ?", dto.GastoID).
		Take(&gasto_map).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, administracion.ErrGastoNotFound
	}

	if err != nil {
		return nil, core.WrapError(err)
	}

	gasto := administracion.Gasto{
		IGasto: administracion.IGasto{
			ID:             gasto_map["id"].(string),
			Descripcion:    nil,
			Monto:          int(gasto_map["monto"].(int64)),
			Proveedor:      gasto_map["proveedor"].(string),
			Moneda:         gasto_map["moneda"].(string),
			Tasa:           int(gasto_map["tasa"].(int64)),
			Fecha:          gasto_map["fecha"].(time.Time),
			Registro:       gasto_map["registro"].(time.Time),
			Actualizacion:  gasto_map["actualizacion"].(time.Time),
			RegistradoPor:  gasto_map["registrado_por"].(string),
			ActualizadoPor: gasto_map["actualizado_por"].(string),
		},
		Total: int((*gasto_map["total"].(*any)).(int64)),
	}

	if gasto_map["descripcion"] != nil {
		desc := gasto_map["descripcion"].(string)
		gasto.Descripcion = &desc
	}

	return &gasto, nil
}
