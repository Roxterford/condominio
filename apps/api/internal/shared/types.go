package shared

import "github.com/Sanaruca/condominio/internal/core/common/quantity"

// AlcanceUnidades representa el impacto de la cuota en el condominio
type AlcanceUnidades struct {
	TotalUnidades int
}

// EjecucionPagos representa el estado financiero de la cuota en tiempo real
type EjecucionPagos struct {
	CantidadPagos  int
	MontoRecaudado quantity.Quantity
}
