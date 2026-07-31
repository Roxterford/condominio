// Deprecated: Modelo legacy de gasto. Usar internal/transacciones/ en su lugar.
package gasto

import (
	"time"

	"github.com/Sanaruca/condominio/internal/core"
	"github.com/Sanaruca/condominio/internal/core/common/audit"
	"github.com/Sanaruca/condominio/internal/core/common/filter"
	"github.com/Sanaruca/condominio/internal/core/common/quantity"
	"github.com/Sanaruca/condominio/internal/core/errors"
	currency "github.com/Sanaruca/condominio/internal/pagos/types/moneda"
)

var (
	ErrGastoNotFound = errors.New(errors.NOT_FOUND, "Gasto no encontrado")
)

type GastoID string

type Gasto struct {
	id          GastoID
	concepto    string
	proveedor   string
	cuota       *string
	monto       quantity.Quantity
	moneda      currency.Moneda
	tasa        quantity.Quantity
	fecha       time.Time
	descripcion *string
	audit       audit.CreationAudit[string]
}

func (g Gasto) ID() GastoID                        { return g.id }
func (g Gasto) Concepto() string                   { return g.concepto }
func (g Gasto) Proveedor() string                  { return g.proveedor }
func (g Gasto) Cuota() *string                     { return g.cuota }
func (g Gasto) Monto() quantity.Quantity           { return g.monto }
func (g Gasto) Moneda() currency.Moneda            { return g.moneda }
func (g Gasto) Tasa() quantity.Quantity            { return g.tasa }
func (g Gasto) Fecha() time.Time                   { return g.fecha }
func (g Gasto) Descripcion() *string               { return g.descripcion }
func (g Gasto) Audit() audit.CreationAudit[string] { return g.audit }

func (g *Gasto) SetCuota(cuotaID string) core.Error {
	if g.cuota != nil {
		return errors.New(errors.CONFLICT, "El gasto '%s' ya tiene una cuota asignada", g.id)
	}
	g.cuota = &cuotaID
	return nil
}

// Monto total en USD
func (g Gasto) Total() quantity.Quantity {

	switch g.moneda {
	case currency.USD:
		return g.monto
	case currency.VED:
		return g.monto.HappyDiv(g.tasa)
	}

	return quantity.Quantity{}
}

func (g Gasto) FilterSpec() filter.Spec {
	return filter.Spec{
		"cuota": filter.TypeString,
	}
}
