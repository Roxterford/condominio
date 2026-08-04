package transaccion

import (
	"github.com/Sanaruca/condominio/internal/core/common/quantity"
	"github.com/Sanaruca/condominio/internal/finanzas/types/roldelmovimiento"
	"github.com/Sanaruca/condominio/internal/finanzas/types/tipodemovimiento"
)

type Movimiento struct {
	id            string
	cuota         *string
	tipo          tipodemovimiento.TipoDeMovimiento
	monto         quantity.Quantity
	rol           roldelmovimiento.RolDelMovimiento
	unidad_codigo *string
	proveedor_id  *string
}

func (m *Movimiento) ID() string                              { return m.id }
func (m *Movimiento) CuotaID() *string                        { return m.cuota }
func (m *Movimiento) Tipo() tipodemovimiento.TipoDeMovimiento { return m.tipo }
func (m *Movimiento) Monto() quantity.Quantity                { return m.monto }
func (m *Movimiento) Rol() roldelmovimiento.RolDelMovimiento  { return m.rol }
func (m *Movimiento) UnidadCodigo() *string                   { return m.unidad_codigo }
func (m *Movimiento) ProveedorID() *string                    { return m.proveedor_id }

func newMovimiento(
	id string,
	tipo tipodemovimiento.TipoDeMovimiento,
	monto quantity.Quantity,
	rol roldelmovimiento.RolDelMovimiento,
	unidad_codigo *string,
	proveedor_id *string,
) *Movimiento {
	return &Movimiento{
		id:            id,
		tipo:          tipo,
		monto:         monto,
		rol:           rol,
		unidad_codigo: unidad_codigo,
		proveedor_id:  proveedor_id,
	}
}

func AssembleMovimiento(
	id string,
	tipo tipodemovimiento.TipoDeMovimiento,
	monto quantity.Quantity,
	rol roldelmovimiento.RolDelMovimiento,
	unidad_codigo *string,
	proveedor_id *string,
) *Movimiento {
	return &Movimiento{
		id:            id,
		tipo:          tipo,
		monto:         monto,
		rol:           rol,
		unidad_codigo: unidad_codigo,
		proveedor_id:  proveedor_id,
	}
}
