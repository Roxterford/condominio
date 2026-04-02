package gasto

import (
	"time"

	"github.com/Sanaruca/condominio/internal/core"
	"github.com/Sanaruca/condominio/internal/core/common/audit"
	currency "github.com/Sanaruca/condominio/internal/pagos/types/moneda"
)

type GastoFactory struct{}

func NewGastoFactory() *GastoFactory {
	return &GastoFactory{}
}

func (f *GastoFactory) Nuevo(
	registrador string,
	concepto string,
	proveedor string,
	monto int,
	moneda currency.Moneda,
	tasa int,
	fecha time.Time,
) (*Gasto, core.Error) {
	return NuevoGasto(registrador, concepto, proveedor, monto, moneda, tasa, fecha)
}

func (f *GastoFactory) Assemble(
	id string,
	concepto string,
	proveedor string,
	cuota *string,
	monto int,
	moneda currency.Moneda,
	tasa int,
	fecha time.Time,
	descripcion *string,
	registrador string,
	creado_en time.Time,
) *Gasto {
	return &Gasto{
		id:          GastoID(id),
		concepto:    concepto,
		proveedor:   proveedor,
		cuota:       cuota,
		monto:       monto,
		moneda:      moneda,
		tasa:        tasa,
		fecha:       fecha,
		descripcion: descripcion,
		audit: audit.CreationAudit[string]{
			CreatedAt: creado_en,
			CreatedBy: registrador,
		},
	}
}
