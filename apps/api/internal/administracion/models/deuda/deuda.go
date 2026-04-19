package deuda

import (
	"time"

	"github.com/Sanaruca/condominio/internal/administracion/models/cuota"
	"github.com/Sanaruca/condominio/internal/core"
	"github.com/Sanaruca/condominio/internal/core/errors"
	"github.com/lucsky/cuid"
)

var (
	ErrDeudaNoEncontrada = errors.New(errors.NOT_FOUND, "Deuda no encontrada")
)

type DeudaFactory struct{}

func NewDeudaFactory() *DeudaFactory {
	return &DeudaFactory{}
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
		monto:    monto,
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
) (*Deuda, core.Error) {
	return &Deuda{
		id:       id,
		cuota:    cuota.CuotaID(cuotaID),
		unidad:   unidad,
		monto:    monto_inicial,
		registro: registro,
		abonos:   abonos,
	}, nil
}

type Deuda struct {
	id       string
	cuota    cuota.CuotaID
	unidad   string
	monto    int // inmutable
	registro time.Time
	abonos   []Abono
}

func (d *Deuda) ID() string             { return d.id }
func (d *Deuda) CuotaID() cuota.CuotaID { return d.cuota }
func (d *Deuda) Unidad() string         { return d.unidad }
func (d *Deuda) Monto() int             { return d.monto }
func (d *Deuda) Registro() time.Time    { return d.registro }
func (d *Deuda) Abonos() []Abono        { return d.abonos }

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
