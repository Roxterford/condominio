package transaccion

import (
	"time"

	"github.com/lucsky/cuid"

	"github.com/Sanaruca/condominio/internal/core"
	"github.com/Sanaruca/condominio/internal/finanzas/models/operacion"
)

type TransaccionFactory struct{}

func NewTransaccionFactory() *TransaccionFactory {
	return &TransaccionFactory{}
}

// Nuevo es el guardian de la integridad del contenedor: exige al menos dos
// operaciones y que todas compartan la misma moneda. Genera identificador y
// marcas de tiempo actuales.
func (f *TransaccionFactory) Nuevo(
	fecha time.Time,
	concepto string,
	registrado_por string,
	operaciones []operacion.Operacion,
) (*Transaccion, core.Error) {
	if len(operaciones) < 2 {
		return nil, ErrOperacionesInsuficientes
	}
	if err := validarMismaMoneda(operaciones); err != nil {
		return nil, err
	}

	return &Transaccion{
		id:             cuid.New(),
		fecha:          fecha,
		concepto:       concepto,
		registrado_por: registrado_por,
		registro:       time.Now().UTC(),
		operaciones:    operaciones,
	}, nil
}

// Assemble reconstruye una transaccion desde la persistencia (infalible).
func (f *TransaccionFactory) Assemble(
	id string,
	fecha time.Time,
	concepto string,
	registrado_por string,
	registro time.Time,
	operaciones []operacion.Operacion,
) *Transaccion {
	return &Transaccion{
		id:             id,
		fecha:          fecha,
		concepto:       concepto,
		registrado_por: registrado_por,
		registro:       registro,
		operaciones:    operaciones,
	}
}

func validarMismaMoneda(operaciones []operacion.Operacion) core.Error {
	moneda := operaciones[0].Moneda()
	for _, op := range operaciones[1:] {
		if op.Moneda() != moneda {
			return ErrMonedasDistintas
		}
	}
	return nil
}
