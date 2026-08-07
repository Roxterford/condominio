package transaccion

import (
	"github.com/Sanaruca/condominio/internal/core/common/quantity"
	"github.com/Sanaruca/condominio/internal/finanzas/types/roldelmovimiento"
	"github.com/Sanaruca/condominio/internal/finanzas/types/tipodemovimiento"
)

type Movimiento interface {
	ID() string
	CuotaID() *string
	Tipo() tipodemovimiento.TipoDeMovimiento
	Monto() quantity.Quantity

	AsACondominio() *MovimientoACondominio
	AsAProveedor() *MovimientoAProveedor
	AsAUnidad() *MovimientoAUnidad
}

type MovimientoBase struct {
	id    string
	cuota *string
	tipo  tipodemovimiento.TipoDeMovimiento
	monto quantity.Quantity
}

func (m *MovimientoBase) ID() string                              { return m.id }
func (m *MovimientoBase) CuotaID() *string                        { return m.cuota }
func (m *MovimientoBase) Tipo() tipodemovimiento.TipoDeMovimiento { return m.tipo }
func (m *MovimientoBase) Monto() quantity.Quantity                { return m.monto }

// Required by `gqlgen` to satisfy the `Movimiento` GraphQL interface.
func (m *MovimientoBase) GetID() string { return m.id }

// Required by `gqlgen` to satisfy the `Movimiento` GraphQL interface.
func (m *MovimientoBase) GetTipo() tipodemovimiento.TipoDeMovimiento { return m.tipo }

// Required by `gqlgen` to satisfy the `Movimiento` GraphQL interface.
func (m *MovimientoBase) GetMonto() float64 { return m.monto.Float() }

// Required by `gqlgen` to satisfy the `Movimiento` GraphQL interface.
func (m *MovimientoBase) GetCuota() *string { return m.cuota }

// Required by `gqlgen` to satisfy the `Movimiento` GraphQL interface.
func (MovimientoBase) IsMovimiento() {}

// Required by `gqlgen` to satisfy the `MovimientoType` GraphQL interface.
func (MovimientoBase) IsMovimientoType() {}

func newMovimiento(
	id string,
	tipo tipodemovimiento.TipoDeMovimiento,
	monto quantity.Quantity,
) *MovimientoBase {
	return &MovimientoBase{
		id:    id,
		tipo:  tipo,
		monto: monto,
	}
}

func AssembleMovimiento(
	id string,
	tipo tipodemovimiento.TipoDeMovimiento,
	monto quantity.Quantity,
	rol roldelmovimiento.RolDelMovimiento,
	unidad_codigo *string,
	proveedor_id *string,
) Movimiento {
	base := MovimientoBase{
		id:    id,
		tipo:  tipo,
		monto: monto,
	}

	switch rol {
	case roldelmovimiento.Condominio:
		return &MovimientoACondominio{
			MovimientoBase: base,
		}
	case roldelmovimiento.Proveedor:
		return &MovimientoAProveedor{
			MovimientoBase: base,
			proveedor:      deref(proveedor_id),
		}
	case roldelmovimiento.Unidad:
		return &MovimientoAUnidad{
			MovimientoBase: base,
			unidad:         deref(unidad_codigo),
		}
	default:
		return &MovimientoACondominio{
			MovimientoBase: base,
		}
	}
}

func deref(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
