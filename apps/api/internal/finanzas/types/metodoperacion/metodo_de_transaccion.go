package metodoperacion

import "github.com/Sanaruca/condominio/internal/core"

type MetodoDeOperacion string

const (
	Efectivo                   MetodoDeOperacion = "EFECTIVO"
	TransferenciaNacional      MetodoDeOperacion = "TRANSFERENCIA_NACIONAL"
	PagoMovil                  MetodoDeOperacion = "PAGO_MOVIL"
	TransferenciaInternacional MetodoDeOperacion = "TRANSFERENCIA_INTERNACIONAL"
	Cheque                     MetodoDeOperacion = "CHEQUE"
	Compensacion               MetodoDeOperacion = "COMPENSACION"
)

func (m MetodoDeOperacion) Validate() core.Error {
	switch m {
	case Efectivo,
		TransferenciaNacional,
		PagoMovil,
		TransferenciaInternacional,
		Cheque,
		Compensacion:
		return nil
	default:
		return core.NewInvalidArgumentError("'%s' no es un metodo de operacion valido", string(m))
	}
}

func (m MetodoDeOperacion) EFECTIVO() bool               { return m == Efectivo }
func (m MetodoDeOperacion) TRANSFERENCIA_NACIONAL() bool { return m == TransferenciaNacional }
func (m MetodoDeOperacion) PAGO_MOVIL() bool             { return m == PagoMovil }
func (m MetodoDeOperacion) TRANSFERENCIA_INTERNACIONAL() bool {
	return m == TransferenciaInternacional
}
func (m MetodoDeOperacion) CHEQUE() bool       { return m == Cheque }
func (m MetodoDeOperacion) COMPENSACION() bool { return m == Compensacion }
