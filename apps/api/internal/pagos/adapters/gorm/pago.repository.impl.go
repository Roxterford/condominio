package gorm

import (
	"context"
	"errors"

	"github.com/Sanaruca/condominio/internal/core"
	gormAdapter "github.com/Sanaruca/condominio/internal/core/adapters/gorm"
	"github.com/Sanaruca/condominio/internal/core/common/filter"
	"github.com/Sanaruca/condominio/internal/pagos/models/pago"

	"gorm.io/gorm"
)

type GORMPagoRepository struct {
	db *gorm.DB
}

func NewGORMPagoRepository(db *gorm.DB) pago.PagoRepository {
	return &GORMPagoRepository{
		db: db,
	}
}

// Count implements [pago.PagoRepository].
func (r *GORMPagoRepository) Count(ctx context.Context, filter filter.Clause) (int, core.Error) {
	count, err := gorm.G[Pago](r.db).Scopes(gormAdapter.GFilter(filter)).Count(ctx, "id")

	if err != nil {
		return 0, core.WrapError(err)
	}

	return int(count), nil

}

// GetByID implements [pago.PagoRepository].
func (r *GORMPagoRepository) GetByID(ctx context.Context, id string) (*pago.Pago, core.Error) {

	dbpago, err := gorm.G[Pago](r.db).Where("id = ?", id).First(ctx)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, pago.ErrPagoNoEncontrado
	}

	if err != nil {
		return nil, core.WrapError(err)
	}

	firma, err := pago.NuevaFirmaFromStore(
		dbpago.Registro,
		dbpago.RegistradoPor,
		dbpago.Actualizacion,
		dbpago.ActualizadoPor,
	)
	if err != nil {
		return nil, core.WrapError(err)
	}

	return pago.NuevoPagoFromStore(
		dbpago.ID,
		dbpago.Unidad,
		dbpago.Fecha,
		dbpago.Metodo,
		dbpago.Monto,
		dbpago.Moneda,
		dbpago.Tasa,
		dbpago.Referencia,
		firma,
	)

}

// Guardar implements [pago.PagoRepository].
func (r *GORMPagoRepository) Guardar(ctx context.Context, pago *pago.Pago) core.Error {
	// Comprobar existencia
	_, err := gorm.G[IPago](r.db).Where("id = ?", pago.ID()).Select("id").Take(ctx)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return r.insertarNuevoPago(ctx, pago)
	} else if err != nil {
		return core.WrapError(err)
	}

	return r.actualizarPagoExistente(ctx, pago)
}

func (r *GORMPagoRepository) insertarNuevoPago(ctx context.Context, pago *pago.Pago) core.Error {
	destinos := r.mapearDestinos(pago.Destinos(), pago.ID())

	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := gorm.G[IPago](tx).Create(ctx, mapToIPago(pago)); err != nil {
			return err
		}

		if len(destinos) > 0 {
			return gorm.G[DestinoDePago](tx).CreateInBatches(ctx, &destinos, 100)
		}

		return nil
	})

	if err != nil {
		return core.WrapError(err)
	}
	return nil
}

func (r *GORMPagoRepository) actualizarPagoExistente(
	ctx context.Context,
	pago *pago.Pago,
) core.Error {
	var destinos_almacenados []string

	err := r.db.Model(&DestinoDePago{}).
		Where("pago = ?", pago.ID()).
		Pluck("id", &destinos_almacenados).Error

	if err != nil {
		return core.WrapError(err)
	}

	nuevos_destinos := pago.ObtenerDiferenciaDeDestinos(destinos_almacenados)
	destinos_de_pago_a_insertar := r.mapearDestinos(nuevos_destinos, pago.ID())

	err = r.db.Transaction(func(tx *gorm.DB) error {
		// Actualizar cabecera del pago
		if _, err := gorm.G[IPago](tx).Where("id = ?", pago.ID()).Updates(ctx, *mapToIPago(pago)); err != nil {
			return err
		}

		// Insertar solo los destinos que faltan
		if len(destinos_de_pago_a_insertar) > 0 {
			return gorm.G[DestinoDePago](tx).CreateInBatches(ctx, &destinos_de_pago_a_insertar, 100)
		}

		return nil
	})

	if err != nil {
		return core.WrapError(err)
	}
	return nil
}
