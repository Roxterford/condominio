package gorm

import (
	"time"

	"github.com/Sanaruca/condominio/internal/administracion/models/deuda"
	"github.com/Sanaruca/condominio/internal/core/common/quantity"
	uadapters "github.com/Sanaruca/condominio/internal/unidades/adapters/gorm"
	"github.com/Sanaruca/condominio/internal/unidades/models/unidad"
)

type DeudaTable struct {
	ID           string
	UnidadID     string `gorm:"column:unidad_id"`
	UnidadCodigo string `gorm:"column:unidad_codigo"`
	Cuota        string
	// Este valor es el monto de la cuota
	Monto int
	// Este valor es el monto restante de la cuota e ira reduciendose a medida que
	// se realicen los pagos
	Deuda         int
	Estado        string
	Registro      time.Time
	Actualizacion time.Time

	Unidad uadapters.Unidad `gorm:"foreignKey:UnidadID;references:ID"`
	Abonos []DestinoDePago  `gorm:"foreignKey:Deuda"`
}

func (DeudaTable) TableName() string { return "deudas" }

func (t DeudaTable) ToDomainDeuda(
	factory *deuda.DeudaFactory,
	qf *quantity.QuantityFactory,
) *deuda.Deuda {

	abonos := make([]deuda.Abono, len(t.Abonos))

	for i, a := range t.Abonos {
		abonos[i] = *factory.AssembleAbono(a.Operacion, qf.Assemble(int64(a.Destinado)), a.Fecha)
	}

	return factory.Assemble(
		t.ID,
		t.Cuota,
		unidad.WrapIDs(t.UnidadID, t.UnidadCodigo),
		t.Monto,
		t.Registro,
		abonos,
	)

}
