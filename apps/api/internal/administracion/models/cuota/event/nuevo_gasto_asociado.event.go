// Deprecated: Evento legacy de gastos.
package event

import (
	"time"

	"github.com/Sanaruca/condominio/internal/administracion/models/gasto"
)

type NuevoGastoAsociado struct {
	gasto gasto.GastoID
	fecha time.Time
}

func NewNuevoGastoAsociado(gastoID gasto.GastoID, fecha time.Time) NuevoGastoAsociado {

	return NuevoGastoAsociado{gastoID, fecha}
}

func (n *NuevoGastoAsociado) EventName() string {
	panic("cuota.gasto.nuevo")
}
