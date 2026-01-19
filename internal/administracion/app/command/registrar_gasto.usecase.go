package command

import (
	"errors"
	"strings"
	"time"

	"github.com/Sanaruca/condominio/internal/administracion"
	"github.com/Sanaruca/condominio/internal/administracion/app/query"
	"github.com/Sanaruca/condominio/internal/core"
	"github.com/Sanaruca/condominio/internal/core/adapters/ozzo"
	"github.com/Sanaruca/condominio/internal/core/context"
	"github.com/Sanaruca/condominio/internal/core/usecase"
	"github.com/Sanaruca/condominio/internal/pagos/types/moneda"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/lucsky/cuid"
	"gorm.io/gorm"
)

type RegistrarGastoDTO struct {
	Proveedor   string
	Monto       int
	Moneda      moneda.Moneda
	Tasa        int
	Fecha       *time.Time
	Descripcion *string
}

type RegistrarGasto usecase.Handler[context.AdminContext, RegistrarGastoDTO, *administracion.Gasto]

func NewRegistrarGasto(obtenerGasto query.ObtenerGasto) RegistrarGasto {
	return &registrarGasto{
		obtenerGasto: obtenerGasto,
	}
}

type registrarGasto struct {
	obtenerGasto query.ObtenerGasto
}

func (uc *registrarGasto) Exec(
	ctx context.AdminContext,
	dto RegistrarGastoDTO,
) (*administracion.Gasto, core.Error) {

	if err := dto.Validate(); err != nil {
		return nil, err
	}

	proveedor, err := gorm.G[administracion.Proveedor](
		ctx.DB,
	).Where("id = ?", dto.Proveedor).
		Select("id").
		First(ctx.Context())

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, administracion.ErrProveedorNoEncontrado
	}

	if err != nil {
		return nil, core.WrapError(err)
	}

	internal_gasto := &administracion.IGasto{
		Id:             cuid.New(),
		Proveedor:      proveedor.ID,
		Monto:          dto.Monto,
		Moneda:         dto.Moneda.String(),
		Tasa:           dto.Tasa,
		Fecha:          time.Now(),
		Descripcion:    dto.Descripcion,
		Registro:       time.Now(),
		Actualizacion:  time.Now(),
		RegistradoPor:  ctx.Session.Usuario.ID,
		ActualizadoPor: ctx.Session.Usuario.ID,
	}

	if dto.Fecha != nil {
		internal_gasto.Fecha = *dto.Fecha
	}

	if err := gorm.G[administracion.IGasto](ctx.DB).Create(ctx.Context(), internal_gasto); err != nil {
		return nil, core.WrapError(err)
	}

	return uc.obtenerGasto.Exec(ctx.BaseContext, query.ObtenerGastoDTO{
		GastoID: internal_gasto.Id,
	})
}

func (dto *RegistrarGastoDTO) Validate() core.Error {

	hoy := time.Now().Truncate(24 * time.Hour)

	if dto.Proveedor != "" {
		dto.Proveedor = strings.TrimSpace(dto.Proveedor)
	}

	if dto.Descripcion != nil {
		*dto.Descripcion = strings.TrimSpace(*dto.Descripcion)
	}

	if dto.Fecha != nil {
		fechaTruncada := dto.Fecha.Truncate(24 * time.Hour)
		dto.Fecha = &fechaTruncada
	} else {
		dto.Fecha = &hoy
	}

	err := validation.ValidateStruct(dto,
		validation.Field(&dto.Proveedor,
			validation.Required.Error("El proveedor es requerido"),
		),
		validation.Field(&dto.Monto,
			validation.Required.Error("El monto es requerido"),
			validation.Min(1).Error("El monto debe ser mayor a 0"),
		),
		validation.Field(&dto.Moneda,
			validation.Required.Error("La moneda es requerida"),
		),
		validation.Field(&dto.Tasa,
			validation.Required.Error("La tasa es requerida"),
			validation.Min(1).Error("La tasa debe ser mayor a 0"),
		),
		validation.Field(&dto.Fecha,
			validation.NilOrNotEmpty.Error("La fecha no puede estar vacía"),
			validation.Max(hoy).Error("La fecha no puede ser futura"),
		),
		validation.Field(&dto.Descripcion,
			validation.NilOrNotEmpty.Error("La descripción no puede estar vacía"),
			validation.Length(0, 500).Error("La descripción no puede exceder los 500 caracteres"),
		),
	)

	if err != nil {
		return ozzo.FirstOzzoErrorAdapter(dto, err)
	}

	// Validar que la moneda tenga un valor válido
	if err := dto.Moneda.Validate(); err != nil {
		return err
	}

	return nil
}
