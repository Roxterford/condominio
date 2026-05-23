package gorm

import (
	"context"
	"errors"
	"log"

	"gorm.io/gorm"

	"github.com/Sanaruca/condominio/internal/core"
	gormAdapter "github.com/Sanaruca/condominio/internal/core/adapters/gorm"
	"github.com/Sanaruca/condominio/internal/core/common"
	"github.com/Sanaruca/condominio/internal/core/common/filter"
	"github.com/Sanaruca/condominio/internal/core/common/quantity"
	"github.com/Sanaruca/condominio/internal/pagos/models/pago"
	"github.com/Sanaruca/condominio/internal/unidades/models/unidad"
)

type GORMPagoRepository struct {
	db *gorm.DB
	qf *quantity.QuantityFactory
	pf *pago.PagoFactory
}

func NewGORMPagoRepository(
	db *gorm.DB,
	quantityFactory *quantity.QuantityFactory,
	pagoFactory *pago.PagoFactory,
) pago.PagoRepository {

	if quantityFactory == nil {
		panic("quantityFactory is nill")
	}
	if pagoFactory == nil {
		panic("pagoFactory is nill")
	}

	return &GORMPagoRepository{
		db: db,
		qf: quantityFactory,
		pf: pagoFactory,
	}
}

// Obtener implements [pago.PagoRepository].
func (r *GORMPagoRepository) Obtener(
	ctx context.Context,
	filter filter.Clause,
	paginator common.Paginator,
) (*common.Paginated[pago.Pago], core.Error) {

	paginator.Sanitize()

	rows, err := gorm.G[Pago](r.db).
		Scopes(
			gormAdapter.GFilter(filter),
			gormAdapter.GPaginate(paginator),
		).
		Find(ctx)

	if err != nil {
		return nil, core.WrapError(err)
	}

	total, err := gorm.G[Pago](r.db).
		Scopes(
			gormAdapter.GFilter(filter),
		).
		Count(ctx, "id")

	if err != nil {
		return nil, core.WrapError(err)
	}

	data := make([]pago.Pago, len(rows))

	for i, p := range rows {
		data[i] = *p.ToDomain(r.pf, r.qf)
	}

	return common.NewPaginated(data, int(total), paginator), nil

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
		unidad.UnidadCodigo(dbpago.Unidad),
		dbpago.Fecha,
		dbpago.Metodo,
		r.qf.Assemble(int64(dbpago.Monto)),
		dbpago.Moneda,
		r.qf.Assemble(int64(dbpago.Tasa)),
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

	log.Printf("\033[36m[DEBUG] insertarNuevoPago: iniciando inserción de pago id=%s con %d destinos\033[0m", pago.ID(), len(destinos))

	err := r.db.Transaction(func(tx *gorm.DB) error {
		log.Printf("\033[36m[DEBUG] insertarNuevoPago: insertando cabecera del pago id=%s\033[0m", pago.ID())
		if err := gorm.G[IPago](tx).Create(ctx, mapToIPago(pago)); err != nil {
			log.Printf("\033[31m[DEBUG] insertarNuevoPago: error al insertar cabecera del pago id=%s: %v\033[0m", pago.ID(), err)
			return err
		}
		log.Printf("\033[32m[DEBUG] insertarNuevoPago: cabecera del pago id=%s insertada correctamente\033[0m", pago.ID())

		if len(destinos) > 0 {
			log.Printf("\033[36m[DEBUG] insertarNuevoPago: insertando %d destinos para pago id=%s\033[0m", len(destinos), pago.ID())
			if err := gorm.G[DestinoDePago](tx).CreateInBatches(ctx, &destinos, 100); err != nil {
				log.Printf("\033[31m[DEBUG] insertarNuevoPago: error al insertar destinos para pago id=%s: %v\033[0m", pago.ID(), err)
				return err
			}
			log.Printf("\033[32m[DEBUG] insertarNuevoPago: destinos para pago id=%s insertados correctamente\033[0m", pago.ID())
		}

		return nil
	})

	if err != nil {
		log.Printf("[DEBUG] insertarNuevoPago: transacción fallida para pago id=%s: %v", pago.ID(), err)
		return core.WrapError(err)
	}

	log.Printf("[DEBUG] insertarNuevoPago: pago id=%s insertado exitosamente", pago.ID())
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
