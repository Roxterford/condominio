package tipoperacion

import (
	"github.com/Sanaruca/condominio/internal/core"
	"github.com/Sanaruca/condominio/internal/core/exception"
)

type TipoDeOperacion string

const (
	Debito  TipoDeOperacion = "DEBITO"
	Credito TipoDeOperacion = "CREDITO"
)

var ErrTipoDeOperacionInvalido = exception.New(
	exception.VALIDATION,
	"Tipo de operacion invalido",
)

func (t TipoDeOperacion) Validate() core.Error {
	switch t {
	case Debito, Credito:
		return nil
	default:
		return ErrTipoDeOperacionInvalido
	}
}

func (t TipoDeOperacion) String() string {
	return string(t)
}
func (t TipoDeOperacion) DEBITO() bool  { return t == Debito }
func (t TipoDeOperacion) CREDITO() bool { return t == Credito }
