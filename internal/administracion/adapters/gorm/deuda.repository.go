package gorm

import (
	"context"
	"errors"

	"github.com/Sanaruca/condominio/internal/administracion"
	"github.com/Sanaruca/condominio/internal/administracion/types/estadodeuda"
	"github.com/Sanaruca/condominio/internal/core"
	"gorm.io/gorm"
)

type DeudaGORMRepository struct {
	db           *gorm.DB
	deudaFactory *administracion.DeudaFactory
}

func NewDeudaGORMRepository(
	db *gorm.DB,
	deudaFactory *administracion.DeudaFactory,
) DeudaGORMRepository {
	return DeudaGORMRepository{db: db, deudaFactory: deudaFactory}
}

func (r *DeudaGORMRepository) GetLastDeudaWhereNotPagada(
	ctx context.Context,
	villa int,
) (*administracion.Deuda, core.Error) {
	deuda, err := gorm.G[Deuda](r.db).Where(
		"villa = ? AND estado <> ?",
		villa,
		estadodeuda.Pagada,
	).Select("id", "deuda").Order("registro asc").Take(ctx)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if err != nil {
		return nil, core.WrapError(err)
	}

	destinos, err := gorm.G[DestinoDePago](r.db).Where("deuda = ?", deuda.ID).Find(ctx)

	if err != nil {
		return nil, core.WrapError(err)
	}

	abonos := make([]administracion.Abono, len(destinos))
	for i, destino := range destinos {
		abonos[i] = *r.deudaFactory.AssembleAbono(destino.Pago, destino.Destinado, destino.Fecha)
	}

	return r.deudaFactory.Assemble(
		deuda.ID,
		deuda.Cuota,
		deuda.Villa,
		deuda.Monto,
		deuda.Registro,
		abonos,
	)
}
