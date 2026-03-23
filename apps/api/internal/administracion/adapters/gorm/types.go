package gorm

import (
	"time"

	"github.com/Sanaruca/condominio/internal/administracion/models/proveedor/tipodeproveedor"
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

func (t Proveedor) TableName() string {
	return "proveedores"
}

type Deuda struct {
	ID    string
	Villa int
	Cuota string
	// Este valor es el monto de la cuota
	Monto int
	// Este valor es el monto restante de la cuota e ira reduciendose a medida que
	// se realicen los pagos
	Deuda         int
	Estado        string
	Registro      time.Time
	Actualizacion time.Time
}

type DestinoDePago struct {
	ID        string
	Pago      string
	Deuda     string
	Destinado int
	Fecha     time.Time
}
