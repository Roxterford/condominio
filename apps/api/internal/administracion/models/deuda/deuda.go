package deuda

import (
	"time"

	"github.com/Sanaruca/condominio/internal/administracion/models/cuota"
	"github.com/Sanaruca/condominio/internal/administracion/models/deuda/estadodeuda"
	"github.com/Sanaruca/condominio/internal/core/common/quantity"
	"github.com/Sanaruca/condominio/internal/core/errors"
)

var (
	ErrDeudaNoEncontrada = errors.New(errors.NOT_FOUND, "Deuda no encontrada")
)

type Deuda struct {
	id       string
	cuota    cuota.CuotaID
	unidad   string
	monto    quantity.Quantity
	registro time.Time
	abonos   []Abono
}

func (d *Deuda) ID() string               { return d.id }
func (d *Deuda) CuotaID() cuota.CuotaID   { return d.cuota }
func (d *Deuda) Unidad() string           { return d.unidad }
func (d *Deuda) Monto() quantity.Quantity { return d.monto }
func (d *Deuda) Registro() time.Time      { return d.registro }
func (d *Deuda) Abonos() []Abono          { return d.abonos }

func (d Deuda) Estado() estadodeuda.EstadoDeDeuda {

	if d.Saldada() {
		return estadodeuda.Saldada
	}

	if d.Restante().Value() == d.monto.Value() {
		return estadodeuda.Pendiente
	}

	return estadodeuda.Abonada
}

func (d *Deuda) Deuda() quantity.Quantity {
	return d.Restante()
}

func (d *Deuda) Saldada() bool {
	return d.Restante().Value() == 0
}

func (d *Deuda) Restante() quantity.Quantity {
	return d.monto.HappySub(d.Abonado())
}

func (d *Deuda) Abonado() quantity.Quantity {
	var abonado quantity.Quantity
	for i, a := range d.abonos {
		if i == 0 {
			abonado = a.monto
			continue
		}
		abonado = abonado.HappyAdd(a.monto)
	}
	return abonado
}
