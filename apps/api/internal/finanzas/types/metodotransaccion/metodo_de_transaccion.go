package metodotransaccion

import "github.com/Sanaruca/condominio/internal/core"

type MetodoDeTransaccion string

const (
	Efectivo                   MetodoDeTransaccion = "EFECTIVO"
	TransferenciaNacional      MetodoDeTransaccion = "TRANSFERENCIA_NACIONAL"
	PagoMovil                  MetodoDeTransaccion = "PAGO_MOVIL"
	TransferenciaInternacional MetodoDeTransaccion = "TRANSFERENCIA_INTERNACIONAL"
	Cheque                     MetodoDeTransaccion = "CHEQUE"
	Compensacion               MetodoDeTransaccion = "COMPENSACION"
)

func (m MetodoDeTransaccion) Validate() core.Error {
	switch m {
	case Efectivo, TransferenciaNacional, PagoMovil, TransferenciaInternacional, Cheque, Compensacion:
		return nil
	default:
		return core.NewInvalidArgumentError("'%s' no es un metodo de transaccion valido", string(m))
	}
}
