package proveedor

import (
	"time"

	"github.com/Sanaruca/condominio/internal/core/common"
	"github.com/Sanaruca/condominio/internal/core/common/filter"
	"github.com/Sanaruca/condominio/internal/core/exception"
)

var (
	ErrProveedorNoEncontrado = exception.New(exception.NOT_FOUND, "Proveedor no encontrado")
	ErrProveedorDuplicado    = exception.New(exception.CONFLICT, "Proveedor ya existe")
)

// Proveedor representa la entidad raíz de nuestro Agregado.
type Proveedor struct {
	id             string
	rif            common.Rif
	nombre         string
	email          common.Email
	telefono       common.Phone
	direccion      *string
	creado_en      time.Time
	actualizado_en time.Time
}

// Getters (Solo lectura, manteniendo la integridad)
func (p *Proveedor) ID() string          { return p.id }
func (p *Proveedor) Nombre() string      { return p.nombre }
func (p *Proveedor) Rif() common.Rif     { return p.rif }
func (p *Proveedor) Email() common.Email { return p.email }
func (p *Proveedor) Telefono() common.Phone {
	return p.telefono
}
func (p *Proveedor) Direccion() *string       { return p.direccion }
func (p *Proveedor) CreadoEn() time.Time      { return p.creado_en }
func (p *Proveedor) ActualizadoEn() time.Time { return p.actualizado_en }

// CambiarDireccion ejemplo de un "Behavioral Method" en DDD
func (p *Proveedor) CambiarDireccion(nuevaDireccion string) {
	p.direccion = &nuevaDireccion
	p.actualizado_en = time.Now()
}

func (Proveedor) FilterSpec() filter.Spec {
	return filter.Spec{
		"id": filter.TypeString,
	}
}
