package administracion

import (
	"time"

	"github.com/Sanaruca/condominio/internal/administracion/types/tipodeproveedor"
	"github.com/Sanaruca/condominio/internal/core/errors"
)

var (
	ErrProveedorNoEncontrado = errors.New(errors.NOT_FOUND, "Proveedor no encontrado")
	ErrorProveedorDuplicado  = errors.New(errors.CONFLICT, "Proveedor ya existe")
)

type Proveedor struct {
	ID            string
	Rif           string
	Nombre        string
	Tipo          tipodeproveedor.TipoDeProveedor
	Email         string
	Telefono      string
	Direccion     *string
	Registro      time.Time
	Actualizacion time.Time
}

func (_ Proveedor) TableName() string {
	return "proveedores"
}
