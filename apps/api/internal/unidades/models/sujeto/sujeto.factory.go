package sujeto

import "github.com/Sanaruca/condominio/internal/core/common"

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
