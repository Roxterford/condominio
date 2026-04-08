package cuota

type CuotaRegular struct {
	CuotaBase
}

func (c *CuotaRegular) AsRegular() *CuotaRegular {
	return c
}

func (c *CuotaRegular) AsEspecial() *CuotaEspecial {
	return nil
}
