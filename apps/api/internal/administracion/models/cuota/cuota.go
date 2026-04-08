package cuota

import (
	"github.com/Sanaruca/condominio/internal/core/common/audit"
	"github.com/Sanaruca/condominio/internal/core/common/filter"
	"github.com/Sanaruca/condominio/internal/core/errors"
)

var (
	ErrCuotaNoEncontrada = errors.New(errors.NOT_FOUND, "Cuota no encontrada")
)

type Cuota interface {
	AsRegular() *CuotaRegular
	AsEspecial() *CuotaEspecial
}

type CuotaID string

type CuotaBase struct {
	ID    CuotaID
	Monto int
	Mes   int
	Anio  int
	Audit audit.FullAudit[string]
}

func (c CuotaBase) FilterSpec() filter.Spec {
	return filter.Spec{
		"id":            filter.TypeString,
		"monto":         filter.TypeInt,
		"mes":           filter.TypeInt,
		"anio":          filter.TypeInt,
		"registro":      filter.TypeUnknown,
		"actualizacion": filter.TypeUnknown,
	}
}
