package gasto

import (
	"time"

	"github.com/Sanaruca/condominio/internal/core"
	"github.com/Sanaruca/condominio/internal/core/common/audit"
	"github.com/Sanaruca/condominio/internal/core/errors"
	currency "github.com/Sanaruca/condominio/internal/pagos/types/moneda"
	"github.com/lucsky/cuid"
)

var (
	ErrGastoNotFound = errors.New(errors.NOT_FOUND, "Gasto no encontrado")
)

type GastoID string

type Gasto struct {
	id          GastoID
	proveedor   string
	cuota       *string
	monto       int
	moneda      currency.Moneda
	tasa        int
	fecha       time.Time
	descripcion *string
	audit       audit.CreationAudit[string]
}

func (g Gasto) ID() GastoID                        { return g.id }
func (g Gasto) Proveedor() string                  { return g.proveedor }
func (g Gasto) Cuota() *string                     { return g.cuota }
func (g Gasto) Monto() int                         { return g.monto }
func (g Gasto) Moneda() currency.Moneda            { return g.moneda }
func (g Gasto) Tasa() int                          { return g.tasa }
func (g Gasto) Fecha() time.Time                   { return g.fecha }
func (g Gasto) Descripcion() *string               { return g.descripcion }
func (g Gasto) Audit() audit.CreationAudit[string] { return g.audit }

func NuevoGasto(
	registrador string, // Usuario que registra el gasto
	proveedor string,
	monto int,
	moneda currency.Moneda,
	tasa int,
	fecha time.Time,
) (*Gasto, core.Error) {

	if proveedor == "" {
		return nil, core.NewValidationError("El proveedor es requerido")
	}
	if monto < 1 {
		return nil, core.NewValidationError("El monto debe ser mayor a 0")
	}
	if moneda != currency.USD && moneda != currency.VED {
		return nil, core.NewValidationError("La moneda debe ser USD o VED")
	}
	if moneda == currency.VED && tasa < 1 {
		return nil, core.NewValidationError("La tasa debe ser mayor a 0 cuando la moneda es VED")
	}
	if fecha.IsZero() {
		fecha = time.Now().UTC()
	}

	return &Gasto{
		id:        GastoID(cuid.New()),
		proveedor: proveedor,
		monto:     monto,
		moneda:    moneda,
		tasa:      tasa,
		fecha:     fecha,
		audit:     audit.NewCreationAudit(registrador),
	}, nil
}

// Monto total en USD
func (g Gasto) Total() int {

	switch g.moneda {
	case currency.USD:
		return g.monto
	case currency.VED:
		return g.monto / g.tasa
	}

	return 0
}
