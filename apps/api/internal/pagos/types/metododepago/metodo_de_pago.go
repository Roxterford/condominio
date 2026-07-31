// Deprecated: Usar internal/transacciones/types/metododetransaccion/ en su lugar.
package metododepago

import "github.com/Sanaruca/condominio/internal/core"

type MetodoDePago string

const (
	Efectivo      MetodoDePago = "EFECTIVO"
	Transferencia MetodoDePago = "TRANSFERENCIA"
	PagoMovil     MetodoDePago = "PAGOMOVIL"
)

func (m MetodoDePago) Validate() core.Error {

	switch m {
	case Efectivo, Transferencia, PagoMovil:
		return nil
	default:
		return core.NewInvalidArgumentError("'%s' no es un metodo de pago valido", string(m))
	}

}
