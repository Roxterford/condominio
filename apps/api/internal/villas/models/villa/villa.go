package villa

import (
	"github.com/Sanaruca/condominio/internal/core/errors"
	"github.com/Sanaruca/condominio/internal/villas/models/villa/estadovilla"
)

var (
	ErrVillaNoEncontrada = errors.New(errors.NOT_FOUND, "Villa no encontrada")
)

type Villa struct {
	id     string
	numero int
	estado estadovilla.EstadoDeVilla
	deuda  int
}

func (v *Villa) ID() string                        { return v.id }
func (v *Villa) Numero() int                       { return v.numero }
func (v *Villa) Estado() estadovilla.EstadoDeVilla { return v.estado }
func (v *Villa) Deuda() int                        { return v.deuda }

func (v *Villa) PoseeDeuda() bool { return v.deuda > 0 }
