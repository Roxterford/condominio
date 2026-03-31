package cache

import (
	"context"
	"time"

	"github.com/Sanaruca/condominio/internal/core/errors"
	"github.com/Sanaruca/condominio/internal/services/tasa"
	"gorm.io/gorm"
)

type TasaCacheRepository interface {
	Obtener(ctx context.Context, tipo tasa.TipoDeCambio, fecha time.Time) (tasa.Tasa, error)
	Guardar(ctx context.Context, tasa tasa.Tasa) error
}

type GormTasaCacheRepository struct {
	db *gorm.DB
}

func NewGormTasaCacheRepository(db *gorm.DB) *GormTasaCacheRepository {
	return &GormTasaCacheRepository{db: db}
}

func (r *GormTasaCacheRepository) Obtener(
	ctx context.Context,
	tipo tasa.TipoDeCambio,
	fecha time.Time,
) (tasa.Tasa, error) {
	var cache TasaCacheModel

	err := r.db.WithContext(ctx).
		Where("tipo = ? AND fecha = ?", tipo, fecha.Truncate(24*time.Hour)).
		First(&cache).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return tasa.Tasa{}, tasa.ErrTasaNoEncontrada
		}
		return tasa.Tasa{}, errors.Wrap(err)
	}

	return tasa.Tasa{
		Valor:  cache.Valor,
		Fuente: cache.Fuente,
		Fecha:  cache.Fecha,
		Tipo:   tasa.TipoDeCambio(cache.Tipo),
		Moneda: cache.Moneda,
	}, nil
}

func (r *GormTasaCacheRepository) Guardar(ctx context.Context, tasa tasa.Tasa) error {
	cache := TasaCacheModel{
		Fecha:  tasa.Fecha.Truncate(24 * time.Hour),
		Tipo:   string(tasa.Tipo),
		Valor:  tasa.Valor,
		Fuente: tasa.Fuente,
		Moneda: tasa.Moneda,
	}

	return r.db.WithContext(ctx).
		Where("tipo = ? AND fecha = ?", tasa.Tipo, cache.Fecha).
		Assign(cache).
		FirstOrCreate(&TasaCacheModel{}).Error
}

type TasaCacheModel struct {
	Fecha  time.Time `gorm:"primaryKey"`
	Tipo   string    `gorm:"primaryKey"`
	Valor  int
	Fuente string
	Moneda string
}

func (TasaCacheModel) TableName() string {
	return "tasas_cache"
}
