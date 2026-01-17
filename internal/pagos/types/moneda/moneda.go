package moneda

import (
	"github.com/Sanaruca/condominio/internal/core"
)

type Moneda string

const (
	USD Moneda = "USD"
	VED Moneda = "VED"
)

func (m Moneda) Validate() core.Error {
	switch m {
	case USD, VED:
		return nil
	default:
		return core.NewInvalidArgumentError("'%s' no es una moneda valida", string(m))
	}
}
