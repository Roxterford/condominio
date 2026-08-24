package transaccion

import (
	"time"

	"github.com/Sanaruca/condominio/internal/core/exception"
	"github.com/Sanaruca/condominio/internal/finanzas/models/operacion"
)

var (
	ErrTransaccionNoEncontrada  = exception.New(exception.NOT_FOUND, "Transaccion no encontrada")
	ErrOperacionesInsuficientes = exception.New(
		exception.INVALID_ARGUMENT,
		"Una transaccion debe contener al menos dos operaciones",
	)
	ErrMonedasDistintas = exception.New(
		exception.INVALID_ARGUMENT,
		"Todas las operaciones de una transaccion deben usar la misma moneda",
	)
)

// Transaccion es el contenedor opcional que agrupa dos o mas operaciones
// relacionadas (p. ej. una compensacion). Solo existe para casos excepcionales:
// una operacion suelta no pertenece a ninguna transaccion.
type Transaccion struct {
	id             string
	fecha          time.Time
	concepto       string
	registrado_por string
	registro       time.Time
	operaciones    []operacion.Operacion
}

func (t *Transaccion) ID() string            { return t.id }
func (t *Transaccion) Fecha() time.Time      { return t.fecha }
func (t *Transaccion) Concepto() string      { return t.concepto }
func (t *Transaccion) RegistradoPor() string { return t.registrado_por }
func (t *Transaccion) Registro() time.Time   { return t.registro }
func (t *Transaccion) Operaciones() []operacion.Operacion {
	return t.operaciones
}
