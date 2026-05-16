package local

import (
	"context"
	"time"

	"github.com/Sanaruca/condominio/internal/core/common/quantity"
	"github.com/Sanaruca/condominio/internal/core/errors"
	"github.com/Sanaruca/condominio/internal/services/tasa"
	"gorm.io/gorm"
)

type GormLocalTasaRepository struct {
	db *gorm.DB
}

func NewGormLocalTasaRepository(db *gorm.DB) *GormLocalTasaRepository {
	return &GormLocalTasaRepository{db: db}
}

func (a *GormLocalTasaRepository) Nombre() string {
	return "local"
}

func (a *GormLocalTasaRepository) ObtenerTasaActual(tipo tasa.TipoDeCambio) (tasa.Tasa, error) {
	return a.obtenerTasaParaFecha(time.Now())
}

func (a *GormLocalTasaRepository) ObtenerTasaPorFecha(
	tipo tasa.TipoDeCambio,
	fecha time.Time,
) (tasa.Tasa, error) {
	return a.obtenerTasaParaFecha(fecha)
}

func (a *GormLocalTasaRepository) obtenerTasaParaFecha(fecha time.Time) (tasa.Tasa, error) {
	ctx := context.Background()

	type Result struct {
		Tasa   int
		Fecha  time.Time
		Moneda string
		Origen string
	}

	var result Result

	err := a.db.WithContext(ctx).
		Model(&TasaDeCambio{}).
		Select("tasa, MAX(fecha) as fecha, moneda, origen").
		Where("fecha <= ?", fecha).
		Group("DATE(fecha), moneda, origen").
		Order("MAX(fecha) DESC").
		Limit(1).
		Scan(&result).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return tasa.Tasa{}, tasa.ErrTasaNoEncontrada
		}
		return tasa.Tasa{}, errors.Wrap(err)
	}

	if result.Tasa == 0 {
		return tasa.Tasa{}, tasa.ErrTasaNoEncontrada
	}

	fuente := "local"
	switch result.Origen {
	case "pagos":
		fuente += " (Pagos Registrados)"
	case "gastos":
		fuente += " (Gastos Registrados)"
	}

	return tasa.Tasa{
		Valor:  quantity.New(int64(result.Tasa), quantity.DEFAULT_SCALE),
		Fuente: fuente,
		Fecha:  result.Fecha,
		Tipo:   tasa.CambioPromedio,
		Moneda: result.Moneda,
	}, nil
}

func (a *GormLocalTasaRepository) ObtenerHistorico(
	tipo tasa.TipoDeCambio,
	desde, hasta time.Time,
) ([]tasa.Tasa, error) {
	ctx := context.Background()

	type Result struct {
		Tasa   int
		Fecha  time.Time
		Moneda string
		Origen string
	}

	var results []Result

	err := a.db.WithContext(ctx).
		Model(&TasaDeCambio{}).
		Select("tasa, fecha, moneda, origen").
		Where("fecha >= ? AND fecha <= ?", desde, hasta).
		Order("fecha DESC").
		Scan(&results).Error

	if err != nil {
		return nil, errors.Wrap(err)
	}

	tasas := make([]tasa.Tasa, 0, len(results))
	for _, r := range results {
		fuente := "local"
		switch r.Origen {
		case "pagos":
			fuente += " (Pagos Registrados)"
		case "gastos":
			fuente += " (Gastos Registrados)"
		}

		tasas = append(tasas, tasa.Tasa{
			Valor:  quantity.New(int64(r.Tasa), quantity.DEFAULT_SCALE),
			Fuente: fuente,
			Fecha:  r.Fecha,
			Tipo:   tipo,
			Moneda: r.Moneda,
		})
	}

	return tasas, nil
}

type TasaDeCambio struct {
	Moneda   string    `gorm:"column:moneda"`
	Tasa     int       `gorm:"column:tasa"`
	Fecha    time.Time `gorm:"column:fecha"`
	Origen   string    `gorm:"column:origen"`
	OrigenID int       `gorm:"column:origen_id"`
}

func (TasaDeCambio) TableName() string {
	return "tasas_de_cambio"
}
