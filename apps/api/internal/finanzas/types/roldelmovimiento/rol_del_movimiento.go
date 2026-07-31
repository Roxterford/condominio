package roldelmovimiento

type RolDelMovimiento string

const (
	Unidad     RolDelMovimiento = "UNIDAD"
	Proveedor  RolDelMovimiento = "PROVEEDOR"
	Condominio RolDelMovimiento = "CONDOMINIO"
)

func (r RolDelMovimiento) Validate() error {
	switch r {
	case Unidad, Proveedor, Condominio:
		return nil
	default:
		return nil
	}
}

func (r RolDelMovimiento) String() string {
	return string(r)
}
