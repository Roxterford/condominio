package roldestinoperacion

type RolDestinoDeOperacion string

const (
	Unidad     RolDestinoDeOperacion = "UNIDAD"
	Proveedor  RolDestinoDeOperacion = "PROVEEDOR"
	Condominio RolDestinoDeOperacion = "CONDOMINIO"
)

func (r RolDestinoDeOperacion) Validate() error {
	switch r {
	case Unidad, Proveedor, Condominio:
		return nil
	default:
		return nil
	}
}

func (r RolDestinoDeOperacion) String() string {
	return string(r)
}

func (r RolDestinoDeOperacion) Unidad() bool     { return r == Unidad }
func (r RolDestinoDeOperacion) Proveedor() bool  { return r == Proveedor }
func (r RolDestinoDeOperacion) Condominio() bool { return r == Condominio }
