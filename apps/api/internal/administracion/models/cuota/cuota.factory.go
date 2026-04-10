package cuota

import (
	"time"

	"github.com/Sanaruca/condominio/internal/administracion/models/cuota/estadoproyecto"
	"github.com/Sanaruca/condominio/internal/core"
	"github.com/Sanaruca/condominio/internal/core/common/audit"
	"github.com/lucsky/cuid"
)

type CuotaFactory struct {
	proyectoFactory *ProyectoFactory
}

func (f *CuotaFactory) ProyectoFactory() *ProyectoFactory {
	return f.proyectoFactory
}

func NewCuotaFactory(proyectoFactory *ProyectoFactory) *CuotaFactory {
	if proyectoFactory == nil {
		panic("proyectoFactory is nil")
	}

	return &CuotaFactory{proyectoFactory: proyectoFactory}
}

func (f *CuotaFactory) NuevaRegular(
	monto int,
	mes int,
	anio int,
	registrador string,
) (*CuotaRegular, core.Error) {

	if monto <= 0 {
		return nil, core.NewValidationError("el monto debe ser mayor a cero")
	}

	if mes < 1 || mes > 12 {
		return nil, core.NewValidationError("el mes debe estar entre 1 y 12")
	}

	if anio < 2000 {
		return nil, core.NewValidationError("el año debe ser mayor a 2000")
	}

	ahora := time.Now()
	id := CuotaID(cuid.New())

	return &CuotaRegular{
		CuotaBase: CuotaBase{
			ID:    id,
			Monto: monto,
			Mes:   mes,
			Anio:  anio,
			Audit: audit.FullAudit[string]{
				CreatedAt: ahora,
				CreatedBy: registrador,
				UpdatedAt: ahora,
				UpdatedBy: registrador,
			},
		},
	}, nil
}

func (f *CuotaFactory) NuevaEspecial(
	monto int,
	titulo string,
	descripcion string,
	justificacion string,
	fecha_limite time.Time,
	interes_por_mora int64,
	registrador string,
) (*CuotaEspecial, core.Error) {

	if monto <= 0 {
		return nil, core.NewValidationError("el monto debe ser mayor a cero")
	}

	_proyecto, err := f.proyectoFactory.Nuevo(
		titulo,
		descripcion,
		justificacion,
		fecha_limite,
		interes_por_mora,
		registrador,
	)
	if err != nil {
		return nil, err
	}

	ahora := time.Now()
	id := CuotaID(cuid.New())

	return &CuotaEspecial{
		CuotaBase: CuotaBase{
			ID:    id,
			Monto: monto,
			Audit: audit.FullAudit[string]{
				CreatedAt: ahora,
				CreatedBy: registrador,
				UpdatedAt: ahora,
				UpdatedBy: registrador,
			},
		},
		Detalles: *_proyecto,
	}, nil
}

func (f *CuotaFactory) AssembleRegular(
	id string,
	monto int,
	mes int,
	anio int,
	creado_en time.Time,
	actualizado_en time.Time,
	registrador string,
) *CuotaRegular {
	return &CuotaRegular{
		CuotaBase: CuotaBase{
			ID:    CuotaID(id),
			Monto: monto,
			Mes:   mes,
			Anio:  anio,
			Audit: audit.FullAudit[string]{
				CreatedAt: creado_en,
				CreatedBy: registrador,
				UpdatedAt: actualizado_en,
				UpdatedBy: registrador,
			},
		},
	}
}

func (f *CuotaFactory) AssembleEspecial(
	id string,
	monto int,
	titulo string,
	descripcion string,
	justificacion string,
	estado string,
	fecha_limite time.Time,
	interes_por_mora int64,
	creado_en time.Time,
	actualizado_en time.Time,
	registrado_por, actualizado_por string,
) *CuotaEspecial {
	_proyecto := f.proyectoFactory.Assemble(
		titulo,
		descripcion,
		justificacion,
		parseEstado(estado),
		fecha_limite,
		interes_por_mora,
		creado_en,
		actualizado_en,
		registrado_por,
		actualizado_por,
	)

	return &CuotaEspecial{
		CuotaBase: CuotaBase{
			ID:    CuotaID(id),
			Monto: monto,
			Audit: audit.FullAudit[string]{
				CreatedAt: creado_en,
				CreatedBy: registrado_por,
				UpdatedAt: actualizado_en,
				UpdatedBy: actualizado_por,
			},
		},
		Detalles: *_proyecto,
	}
}

func parseEstado(s string) estadoproyecto.EstadoDeProyecto {
	switch s {
	case "BORRADOR":
		return estadoproyecto.BORRADOR
	case "PUBLICADO":
		return estadoproyecto.ACTIVO
	case "CERRADO":
		return estadoproyecto.CERRADO
	default:
		return estadoproyecto.BORRADOR
	}
}
