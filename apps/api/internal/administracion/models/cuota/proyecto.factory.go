package cuota

import (
	"time"

	"github.com/Sanaruca/condominio/internal/administracion/models/cuota/estadoproyecto"
	"github.com/Sanaruca/condominio/internal/core"
	"github.com/Sanaruca/condominio/internal/core/common"
	"github.com/Sanaruca/condominio/internal/core/common/audit"
)

type ProyectoFactory struct{}

func NewProyectoFactory() *ProyectoFactory {
	return &ProyectoFactory{}
}

func (f *ProyectoFactory) Nuevo(
	titulo string,
	descripcion string,
	justificacion string,
	fecha_limite time.Time,
	interes_por_mora int64,
	registrador string,
) (*Proyecto, core.Error) {

	if descripcion == "" {
		return nil, core.NewValidationError("la descripción no puede estar vacía")
	}

	ahora := time.Now()
	return &Proyecto{
		titulo:           titulo,
		estado:           estadoproyecto.BORRADOR,
		descripcion:      descripcion,
		justificacion:    justificacion,
		fecha_limite:     fecha_limite,
		interes_por_mora: common.NewPercentage(interes_por_mora, 4),
		Audit: audit.FullAudit[string]{
			CreatedAt: ahora,
			CreatedBy: registrador,
			UpdatedAt: ahora,
			UpdatedBy: registrador,
		},
	}, nil
}

func (f *ProyectoFactory) Assemble(
	titulo string,
	descripcion string,
	justificacion string,
	estado estadoproyecto.EstadoDeProyecto,
	fecha_limite time.Time,
	interes_por_mora int64,
	creado_en time.Time,
	actualizado_en time.Time,
	registrado_por, actualizado_por string,

) *Proyecto {
	return &Proyecto{
		titulo:           titulo,
		estado:           estado,
		descripcion:      descripcion,
		justificacion:    justificacion,
		fecha_limite:     fecha_limite,
		interes_por_mora: common.NewPercentage(interes_por_mora, 4),
		Audit: audit.FullAudit[string]{
			CreatedAt: creado_en,
			CreatedBy: registrado_por,
			UpdatedAt: actualizado_en,
			UpdatedBy: actualizado_por,
		},
	}
}
