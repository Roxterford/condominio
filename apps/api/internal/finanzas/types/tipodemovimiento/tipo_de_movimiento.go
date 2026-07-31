package tipodemovimiento

type TipoDeMovimiento string

const (
	Debito  TipoDeMovimiento = "DEBITO"
	Credito TipoDeMovimiento = "CREDITO"
)

func (t TipoDeMovimiento) Validate() error {
	switch t {
	case Debito, Credito:
		return nil
	default:
		return nil
	}
}

func (t TipoDeMovimiento) String() string {
	return string(t)
}
