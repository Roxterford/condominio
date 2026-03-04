package villas

import (
	"time"

	"github.com/Sanaruca/condominio/internal/core/errors"
)

var (
	ErrDeudaNoEncontrada = errors.New(errors.NOT_FOUND, "Deuda no encontrada")
	ErrVillaNoEncontrada = errors.New(errors.NOT_FOUND, "Villa no encontrada")
)

type Deuda struct {
	id       string
	cuotaID  string
	villa    int
	monto    int // inmutable
	registro time.Time
	abonos   []Abono
}

func (d *Deuda) ID() string          { return d.id }
func (d *Deuda) CuotaID() string     { return d.cuotaID }
func (d *Deuda) Villa() int          { return d.villa }
func (d *Deuda) Monto() int          { return d.monto }
func (d *Deuda) Registro() time.Time { return d.registro }
func (d *Deuda) Abonos() []Abono     { return d.abonos }

func (d *Deuda) Saldada() bool {
	return d.Restante() == 0
}

func (d *Deuda) Restante() int {
	return d.monto - d.Abonado()
}

func (d *Deuda) Abonado() int {
	var abonado int
	for _, a := range d.abonos {
		abonado += a.monto
	}
	return abonado
}

type Abono struct {
	pagoID string
	monto  int
	fecha  time.Time
}

func (a *Abono) PagoID() string   { return a.pagoID }
func (a *Abono) Monto() int       { return a.monto }
func (a *Abono) Fecha() time.Time { return a.fecha }
