package tipodeproveedor

import "github.com/Sanaruca/condominio/internal/core"

type TipoDeProveedor string

const (
	PersonaNatural TipoDeProveedor = "PERSONA_NATURAL"
	Compania       TipoDeProveedor = "COMPANIA"
)

func (t TipoDeProveedor) String() string {
	return string(t)
}

func (t TipoDeProveedor) Validate() core.Error {
	switch t {
	case PersonaNatural, Compania:
		return nil
	default:
		return core.NewInvalidArgumentError("'%s' no es un tipo de proveedor valido", string(t))
	}
}
