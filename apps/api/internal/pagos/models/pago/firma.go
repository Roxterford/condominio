// Deprecated: Modelo legacy de pagos.
package pago

import (
	"time"

	"github.com/Sanaruca/condominio/internal/core"
)

var (
	ErrFirmaFechaRegistroFutura = core.NewInvalidArgumentError(
		"Error de consistencia, la fecha de registro no puede ser futura",
	)
	ErrFirmaFechaActualizacionFutura = core.NewInvalidArgumentError(
		"Error de consistencia, la fecha de actualización no puede ser futura",
	)
	ErrFirmaUsuarioNulo = core.NewInvalidArgumentError(
		"El usuario no puede ser nulo",
	)
)

type Firma struct {
	registradoEn   time.Time
	registradoPor  string
	actualizadoEn  time.Time
	actualizadoPor string
}

func (f Firma) ActualizadoPor() string   { return f.actualizadoPor }
func (f Firma) ActualizadoEn() time.Time { return f.actualizadoEn }
func (f Firma) RegistradoPor() string    { return f.registradoPor }
func (f Firma) RegistradoEn() time.Time  { return f.registradoEn }

func NuevaFirma(usuarioID string) (*Firma, core.Error) {

	ahora := time.Now()
	firma := &Firma{
		registradoEn:   ahora,
		registradoPor:  usuarioID,
		actualizadoEn:  ahora,
		actualizadoPor: usuarioID,
	}

	return firma, firma.Validate()
}

func NuevaFirmaFromStore(
	registradoEn time.Time,
	registradoPor string,
	actualizadoEn time.Time,
	actualizadoPor string,
) (Firma, core.Error) {
	if registradoEn.After(time.Now()) {
		return Firma{}, ErrFirmaFechaRegistroFutura
	}
	if actualizadoEn.After(time.Now()) {
		return Firma{}, ErrFirmaFechaActualizacionFutura
	}

	return Firma{
		registradoEn:   registradoEn,
		registradoPor:  registradoPor,
		actualizadoEn:  actualizadoEn,
		actualizadoPor: actualizadoPor,
	}, nil

}

// Actualizar crea una nueva firma basada en la anterior (Inmutabilidad)
func (f Firma) Actualizar(usuarioID string) Firma {
	return Firma{
		registradoEn:   f.registradoEn,
		registradoPor:  f.registradoPor,
		actualizadoEn:  time.Now(),
		actualizadoPor: usuarioID,
	}
}

func (f Firma) Validate() core.Error {
	if f.registradoPor == "" {
		return ErrFirmaUsuarioNulo
	}

	if f.actualizadoPor == "" {
		return ErrFirmaUsuarioNulo
	}

	if f.registradoEn.After(time.Now()) {
		return ErrFirmaFechaRegistroFutura
	}
	if f.actualizadoEn.After(time.Now()) {
		return ErrFirmaFechaActualizacionFutura
	}
	return nil
}
