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
	ID() CuotaID
	Monto() int
	Mes() int
	Anio() int
	AsRegular() *CuotaRegular
	AsEspecial() *CuotaEspecial
}

type CuotaID string

func (id CuotaID) String() string { return string(id) }

type CuotaBase struct {
	id    CuotaID
	monto int
	mes   int
	anio  int
	Audit audit.FullAudit[string]
}

func (c CuotaBase) ID() CuotaID { return c.id }
func (c CuotaBase) Monto() int  { return c.monto }
func (c CuotaBase) Mes() int    { return c.mes }
func (c CuotaBase) Anio() int   { return c.anio }

func (c *CuotaBase) SetID(id CuotaID)                       { c.id = id }
func (c *CuotaBase) SetMonto(monto int)                     { c.monto = monto }
func (c *CuotaBase) SetMes(mes int)                         { c.mes = mes }
func (c *CuotaBase) SetAnio(anio int)                       { c.anio = anio }
func (c *CuotaBase) SetAudit(audit audit.FullAudit[string]) { c.Audit = audit }

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
