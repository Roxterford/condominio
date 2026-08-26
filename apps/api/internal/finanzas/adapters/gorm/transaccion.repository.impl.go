package gorm

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/Sanaruca/condominio/internal/core"
	"github.com/Sanaruca/condominio/internal/core/common/quantity"
	"github.com/Sanaruca/condominio/internal/finanzas/models/operacion"
	"github.com/Sanaruca/condominio/internal/finanzas/models/transaccion"
)

type GORMTransaccionRepository struct {
	db *gorm.DB
	f  *transaccion.TransaccionFactory
	of *operacion.OperacionFactory
	qf *quantity.QuantityFactory
}

func NewGORMTransaccionRepository(
	db *gorm.DB,
	transaccionFactory *transaccion.TransaccionFactory,
	operacionFactory *operacion.OperacionFactory,
	quantityFactory *quantity.QuantityFactory,
) transaccion.TransaccionRepository {
	if db == nil {
		panic("db is nil")
	}
	if transaccionFactory == nil {
		panic("transaccionFactory is nil")
	}
	if operacionFactory == nil {
		panic("operacionFactory is nil")
	}
	if quantityFactory == nil {
		panic("quantityFactory is nil")
	}

	return &GORMTransaccionRepository{
		db: db,
		f:  transaccionFactory,
		of: operacionFactory,
		qf: quantityFactory,
	}
}

// Guardar persiste el contenedor junto al puente con sus operaciones. Es
// inmutable: si el id ya existe, la escritura es un no-op.
func (r *GORMTransaccionRepository) Guardar(
	ctx context.Context,
	tx *transaccion.Transaccion,
) core.Error {
	_, err := gorm.G[Transaccion](r.db).Where("id = ?", tx.ID()).Select("id").Take(ctx)

	if err == nil {
		return nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return core.WrapError(err)
	}

	transaccionTable := Transaccion{
		ID:            tx.ID(),
		Fecha:         tx.Fecha(),
		Concepto:      tx.Concepto(),
		RegistradoPor: tx.RegistradoPor(),
		Registro:      tx.Registro(),
	}

	pasarela := make([]TransaccionOperacion, len(tx.Operaciones()))
	for i, op := range tx.Operaciones() {
		pasarela[i] = TransaccionOperacion{
			ID:            tx.ID() + "-" + op.ID(),
			TransaccionID: tx.ID(),
			OperacionID:   op.ID(),
			Posicion:      i,
		}
	}

	err = r.db.Transaction(func(dbTX *gorm.DB) error {
		if err := gorm.G[Transaccion](dbTX).Create(ctx, &transaccionTable); err != nil {
			return err
		}
		if err := gorm.G[TransaccionOperacion](dbTX).CreateInBatches(ctx, &pasarela, 100); err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return core.WrapError(err)
	}

	return nil
}

func (r *GORMTransaccionRepository) ObtenerPorID(
	ctx context.Context,
	id string,
) (*transaccion.Transaccion, core.Error) {
	row, err := gorm.G[Transaccion](r.db).Where("id = ?", id).First(ctx)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, transaccion.ErrTransaccionNoEncontrada
	}
	if err != nil {
		return nil, core.WrapError(err)
	}

	puentes, err := gorm.G[TransaccionOperacion](r.db).
		Where("transaccion_id = ?", id).
		Order("posicion asc").
		Find(ctx)
	if err != nil {
		return nil, core.WrapError(err)
	}

	operaciones := make([]operacion.Operacion, len(puentes))
	for i, puente := range puentes {
		fila, err := gorm.G[Operacion](r.db).Where("id = ?", puente.OperacionID).First(ctx)
		if err != nil {
			return nil, core.WrapError(err)
		}
		operaciones[i] = *toDomainOperacion(fila, r.qf, r.of)
	}

	return r.f.Assemble(
		row.ID,
		row.Fecha,
		row.Concepto,
		row.RegistradoPor,
		row.Registro,
		operaciones,
	), nil
}
