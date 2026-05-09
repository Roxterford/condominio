package deuda

import (
	"time"

	"github.com/lucsky/cuid"

	"github.com/Sanaruca/condominio/internal/administracion/models/cuota"
	"github.com/Sanaruca/condominio/internal/core"
	"github.com/Sanaruca/condominio/internal/core/common/quantity"
)

type DeudaFactory struct {
	qf *quantity.QuantityFactory
}

func NewDeudaFactory(QuantityFactory *quantity.QuantityFactory) *DeudaFactory {

	if QuantityFactory == nil {
		panic("qf is nil")
	}

	return &DeudaFactory{QuantityFactory}
}

func (f DeudaFactory) NuevaDeuda(
	cuotaID cuota.CuotaID,
	unidad string,
	monto int,
) (*Deuda, core.Error) {
	if unidad == "" {
		return nil, core.NewValidationError("la unidad es requerida")
	}
	if monto < 1 {
		return nil, core.NewValidationError("el monto debe ser mayor a cero")
	}

	return &Deuda{
		id:       cuid.New(),
		cuota:    cuotaID,
		unidad:   unidad,
		monto:    f.qf.Assemble(int64(monto)),
		registro: time.Now().UTC(),
		abonos:   []Abono{},
	}, nil
}

func (f DeudaFactory) Assemble(
	id string,
	cuotaID string,
	unidad string,
	monto_inicial int,
	registro time.Time,
	abonos []Abono,
) *Deuda {
	return &Deuda{
		id:       id,
		cuota:    cuota.CuotaID(cuotaID),
		unidad:   unidad,
		monto:    f.qf.Assemble(int64(monto_inicial)),
		registro: registro,
		abonos:   abonos,
	}
}
