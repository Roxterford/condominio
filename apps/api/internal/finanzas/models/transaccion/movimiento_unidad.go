package transaccion

type MovimientoAUnidad struct {
	MovimientoBase
	unidad string
}

func (m *MovimientoAUnidad) Unidad() string {
	return m.unidad
}

func (m *MovimientoAUnidad) AsACondominio() *MovimientoACondominio {
	return nil
}

func (m *MovimientoAUnidad) AsAProveedor() *MovimientoAProveedor {
	return nil
}

func (m *MovimientoAUnidad) AsAUnidad() *MovimientoAUnidad {
	return m
}
