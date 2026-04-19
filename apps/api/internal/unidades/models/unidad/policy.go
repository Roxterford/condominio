package unidad

import "github.com/Sanaruca/condominio/internal/unidades/models/unidad/estadounidad"

// FacturacionPolicy define la lógica para determinar qué se cobra.
// Esto permite que cada condominio tenga sus propias reglas.
type FacturacionPolicy interface {
	GeneraDeudaPorCuotaRegular(estado estadounidad.EstadoDeUnidad) bool
	GeneraDeudaPorCuotaEspecial(estado estadounidad.EstadoDeUnidad) bool
}
