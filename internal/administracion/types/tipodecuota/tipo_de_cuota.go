package tipodecuota

import "github.com/Sanaruca/condominio/internal/core"

type TipoDeCuota string

const (
	Regular  TipoDeCuota = "REGULAR"
	Especial TipoDeCuota = "ESPECIAL"
	Semilla  TipoDeCuota = "SEMILLA"
)

func (t TipoDeCuota) Validate() core.Error {

	switch t {
	case Regular, Especial, Semilla:
		return nil
	}

	return core.NewInvalidArgumentError("'%s' no es un tipo de cuota valido", string(t))
}
