package gorm

import (
	"time"

	"github.com/Sanaruca/condominio/internal/administracion/models/gasto"
	"github.com/Sanaruca/condominio/internal/administracion/models/proveedor"
	"github.com/Sanaruca/condominio/internal/pagos/types/moneda"
)

type Proveedor struct {
	ID            string
	Rif           string
	Nombre        string
	Email         string
	Telefono      string
	Direccion     *string
	Registro      time.Time
	Actualizacion time.Time
}

func (t Proveedor) TableName() string {
	return "proveedores"
}

func (t Proveedor) ToDomainProveedor(factory *proveedor.ProveedorFactory) *proveedor.Proveedor {
	return factory.Assemble(
		t.ID,
		t.Rif,
		t.Nombre,
		t.Email,
		t.Telefono,
		t.Direccion,
		t.Registro,
		t.Actualizacion,
	)

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

type Gasto struct {
	ID             string
	Concepto       string
	Proveedor      string
	Cuota          *string
	Monto          int
	Moneda         moneda.Moneda
	Tasa           int
	Fecha          time.Time
	Descripcion    *string
	Registro       time.Time
	Registrado_por string
}

func (t Gasto) ToDomainGasto(factory *gasto.GastoFactory) *gasto.Gasto {

	return factory.Assemble(
		t.ID,
		t.Concepto,
		t.Proveedor,
		t.Cuota,
		t.Monto,
		t.Moneda,
		t.Tasa,
		t.Fecha,
		t.Descripcion,
		t.Registrado_por,
		t.Registro,
	)

}

func (Gasto) TableName() string {
	return "internal_gastos"
}

type GastoView struct {
	Gasto
	Total int
}

func (t GastoView) TableName() string {
	return "gastos"
}
