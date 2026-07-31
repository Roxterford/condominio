// Deprecated: Factory legacy de gasto. Usar internal/transacciones/TransaccionFactory en su lugar.
package gasto

import (
	"time"

	"github.com/lucsky/cuid"

	"github.com/Sanaruca/condominio/internal/core"
	"github.com/Sanaruca/condominio/internal/core/common/audit"
	"github.com/Sanaruca/condominio/internal/core/common/quantity"
	currency "github.com/Sanaruca/condominio/internal/pagos/types/moneda"
)

type GastoFactory struct {
}

func NewGastoFactory() *GastoFactory {

	return &GastoFactory{}
}

func (f *GastoFactory) Nuevo(
	registrador string,
	concepto string,
	proveedor string,
	monto quantity.Quantity,
	moneda currency.Moneda,
	tasa quantity.Quantity,
	fecha time.Time,
) (*Gasto, core.Error) {

	if concepto == "" {
		return nil, core.NewValidationError("El concepto es requerido")
	}

	if proveedor == "" {
		return nil, core.NewValidationError("El proveedor es requerido")
	}
	if monto.Value() < 1 {
		return nil, core.NewValidationError("El monto debe ser mayor a 0")
	}
	if moneda != currency.USD && moneda != currency.VED {
		return nil, core.NewValidationError("La moneda debe ser USD o VED")
	}
	if moneda == currency.VED && tasa.Value() < 1 {
		return nil, core.NewValidationError("La tasa debe ser mayor a 0 cuando la moneda es VED")
	}
	if fecha.IsZero() {
		fecha = time.Now().UTC()
	}

	return &Gasto{
		id:        GastoID(cuid.New()),
		concepto:  concepto,
		proveedor: proveedor,
		monto:     monto,
		moneda:    moneda,
		tasa:      tasa,
		fecha:     fecha,
		audit:     audit.NewCreationAudit(registrador),

		cuota:       new(string),
		descripcion: new(string),
	}, nil

}

func (f *GastoFactory) Assemble(
	id string,
	concepto string,
	proveedor string,
	cuota *string,
	monto quantity.Quantity,
	moneda currency.Moneda,
	tasa quantity.Quantity,
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
