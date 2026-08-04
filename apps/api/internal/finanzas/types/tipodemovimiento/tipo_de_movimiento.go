package tipodemovimiento

import (
	"github.com/Sanaruca/condominio/internal/core"
	"github.com/Sanaruca/condominio/internal/core/errors"
)

type TipoDeMovimiento string

const (
	Debito  TipoDeMovimiento = "DEBITO"
	Credito TipoDeMovimiento = "CREDITO"
)

var ErrTipoDeMovimientoInvalido = errors.New(
	errors.VALIDATION,
	"Tipo de movimiento invalido",
)

func (t TipoDeMovimiento) Validate() core.Error {
	switch t {
	case Debito, Credito:
		return nil
	default:
		return ErrTipoDeMovimientoInvalido
	}
}

func (t TipoDeMovimiento) String() string {
	return string(t)
}
func (t TipoDeMovimiento) DEBITO() bool  { return t == Debito }
func (t TipoDeMovimiento) CREDITO() bool { return t == Credito }
