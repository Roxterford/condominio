package transaccion

type MovimientoAProveedor struct {
	MovimientoBase
	proveedor string
}

func (m *MovimientoAProveedor) Proveedor() string {
	return m.proveedor
}

func (m *MovimientoAProveedor) AsACondominio() *MovimientoACondominio {
	return nil
}

func (m *MovimientoAProveedor) AsAProveedor() *MovimientoAProveedor {
	return m
}

func (m *MovimientoAProveedor) AsAUnidad() *MovimientoAUnidad {
	return nil
}

// Required by `gqlgen` to satisfy the `Movimiento` GraphQL interface.
func (MovimientoAProveedor) IsMovimiento() {}

// Required by `gqlgen` to satisfy the `MovimientoType` GraphQL interface.
func (MovimientoAProveedor) IsMovimientoType() {}
