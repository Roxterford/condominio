package villa

import "github.com/Sanaruca/condominio/internal/villas/models/villa/estadovilla"

// FacturacionPolicy define la lógica para determinar qué se cobra.
// Esto permite que cada condominio tenga sus propias reglas.
type FacturacionPolicy interface {
	GeneraDeudaPorCuotaRegular(estado estadovilla.EstadoDeVilla) bool
	GeneraDeudaPorCuotaEspecial(estado estadovilla.EstadoDeVilla) bool
}
