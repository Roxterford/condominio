package sujeto

import (
	"time"

	"github.com/lucsky/cuid"

	"github.com/Sanaruca/condominio/internal/core"
	"github.com/Sanaruca/condominio/internal/core/common"
	"github.com/Sanaruca/condominio/internal/core/common/audit"
)

type SujetoFactory struct {
	ef *common.EmailFactory
	ff *common.PhoneFactory
}

func NewSujetoFactory(
	emailFactory *common.EmailFactory,
	phoneFactory *common.PhoneFactory,
) *SujetoFactory {
	return &SujetoFactory{emailFactory, phoneFactory}
}

func (f SujetoFactory) Nuevo(
	tipo TipoDeSujeto,
	documento string,
	nombres, apellidos, razonSocial, representanteID, email, telefono string,
	registrador string,
) (Sujeto, core.Error) {

	if documento == "" {
		return nil, core.NewValidationError("el documento de identidad es requerido")
	}
	if email == "" {
		return nil, core.NewValidationError("el email es requerido")
	}
	if telefono == "" {
		return nil, core.NewValidationError("el telefono es requerido")
	}

	ahora := time.Now()
	id := SujetoID(cuid.New())
	auditoria := audit.FullAudit[string]{
		CreatedAt: ahora,
		CreatedBy: registrador,
		UpdatedAt: ahora,
		UpdatedBy: registrador,
	}

	base := sujeto{
		id:       id,
		cedula:   common.AssembleRif(documento),
		email:    f.ef.Assemble(email),
		telefono: f.ff.Assemble(telefono),
		audit:    auditoria,
	}

	switch tipo {
	case PERSONA_NATURAL:
		return &Persona{
			sujeto:   base,
			nombres:   nombres,
			apellidos: apellidos,
		}, nil
	case ENTE_JURIDICO:
		return &Ente{
			sujeto:       base,
			razon_social: razonSocial,
		}, nil
	default:
		return nil, core.NewValidationError("tipo de sujeto invalido")
	}
}

func (f SujetoFactory) AssembleEnte(
	id,
	rif,
	razon_social,
	email,
	telefono string,
	representante Persona,
) *Ente {

	return &Ente{
		sujeto: sujeto{
			id:       SujetoID(id),
			cedula:   common.AssembleRif(rif),
			email:    f.ef.Assemble(email),
			telefono: f.ff.Assemble(telefono),
		},
		representante: representante,
		razon_social:  razon_social,
	}

}

func (f SujetoFactory) AssemblePersona(
	id,
	cedula,
	nombres,
	apellidos,
	email,
	telefono string,
) *Persona {

	return &Persona{
		sujeto: sujeto{
			id:       SujetoID(id),
			cedula:   common.AssembleRif(cedula),
			email:    f.ef.Assemble(email),
			telefono: f.ff.Assemble(telefono),
		},
		nombres:   nombres,
		apellidos: apellidos,
	}

}
