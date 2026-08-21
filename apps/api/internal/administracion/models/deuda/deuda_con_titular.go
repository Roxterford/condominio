package deuda

import "github.com/Sanaruca/condominio/internal/unidades/models/sujeto"

// DeudaConTitular es un read-model que combina una deuda con el titular
// primario de su unidad. Se construye en la capa de persistencia con una
// sola consulta y no es un agregado del dominio.
type DeudaConTitular struct {
	Deuda
	Titular sujeto.Titular
}
