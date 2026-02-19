package administracion

import (
	"time"

	"github.com/Sanaruca/condominio/internal/core/common/filter"
	"github.com/Sanaruca/condominio/internal/core/errors"
)

var (
	ErrCuotaNoEncontrada = errors.New(errors.NOT_FOUND, "Cuota no encontrada")
)

type Cuota struct {
	ID            string
	Monto         int
	Mes           int
	Anio          int
	Registro      time.Time
	Actualizacion time.Time
}

func (c Cuota) TableName() string {
	return "cuotas"
}

func (c Cuota) FilterSpec() filter.Spec {
	return filter.Spec{"id": filter.String, "monto": filter.Int, "mes": filter.Int, "anio": filter.Int, "registro": filter.Unknown, "actualizacion": filter.Unknown}
}
