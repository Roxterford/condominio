package pago

import (
	"time"

	"github.com/lucsky/cuid"

	"github.com/Sanaruca/condominio/internal/core/common/quantity"
	"github.com/Sanaruca/condominio/internal/core/errors"
)

type Destino struct {
	id        string
	deuda     string
	destinado quantity.Quantity
	fecha     time.Time
	registro  time.Time
}

func (d *Destino) ID() string                   { return d.id }
func (d *Destino) Deuda() string                { return d.deuda }
func (d *Destino) Destinado() quantity.Quantity { return d.destinado }
func (d *Destino) Fecha() time.Time             { return d.fecha }
func (d *Destino) Regsitro() time.Time          { return d.registro }

func NuevoDestino(deuda string, destinado quantity.Quantity) (*Destino, errors.CoreError) {

	id := cuid.New()
	if deuda == "" {
		return nil, errors.New(
			errors.INVALID_ARGUMENT,
			"La deuda no puede ser nula",
		)
	}
	if destinado.Value() <= 0 {
		return nil, errors.New(
			errors.INVALID_ARGUMENT,
			"El destinado no puede ser menor o igual a cero",
		)
	}

	return &Destino{
		id:        id,
		deuda:     deuda,
		destinado: destinado,
		fecha:     time.Now(),
	}, nil
}
