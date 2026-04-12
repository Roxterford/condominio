package shared

import "github.com/Sanaruca/condominio/internal/core/common/quantity"

// AlcanceVillas representa el impacto de la cuota en el condominio
type AlcanceVillas struct {
	TotalVillas int
}

// EjecucionPagos representa el estado financiero de la cuota en tiempo real
type EjecucionPagos struct {
	CantidadPagos  int
	MontoRecaudado quantity.Quantity
}
