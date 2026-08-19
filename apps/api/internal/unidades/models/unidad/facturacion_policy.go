package unidad

import (
	"github.com/Sanaruca/condominio/internal/unidades/models/unidad/estadounidad"
)

// FacturacionPolicyPorDefecto implementa la política estándar según la LPH:
// generan deuda las unidades Activas, Inhabitadas y EnLitigio; no generan
// las Exentas, Suspendidas ni en Preventa.
type FacturacionPolicyPorDefecto struct{}

func NewFacturacionPolicyPorDefecto() FacturacionPolicyPorDefecto {
	return FacturacionPolicyPorDefecto{}
}

func (FacturacionPolicyPorDefecto) GeneraDeudaPorCuotaRegular(
	estado estadounidad.EstadoDeUnidad,
) bool {
	return generaDeuda(estado)
}

func (FacturacionPolicyPorDefecto) GeneraDeudaPorCuotaEspecial(
	estado estadounidad.EstadoDeUnidad,
) bool {
	return generaDeuda(estado)
}

func generaDeuda(estado estadounidad.EstadoDeUnidad) bool {
	switch estado {
	case estadounidad.Activa, estadounidad.Inhabitada, estadounidad.EnLitigio:
		return true
	default:
		return false
	}
}
