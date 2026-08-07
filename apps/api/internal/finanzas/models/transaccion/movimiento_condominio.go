package transaccion

type MovimientoACondominio struct {
	MovimientoBase
}

func (m *MovimientoACondominio) AsACondominio() *MovimientoACondominio {
	return m
}

func (m *MovimientoACondominio) AsAProveedor() *MovimientoAProveedor {
	return nil
}

func (m *MovimientoACondominio) AsAUnidad() *MovimientoAUnidad {
	return nil
}
