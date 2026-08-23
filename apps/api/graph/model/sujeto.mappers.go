package model

import (
	"github.com/Sanaruca/condominio/internal/unidades/models/sujeto"
)

func TitularFromDomain(titular sujeto.Titular) Titular {
	if titular == nil {
		return nil
	}

	ente := titular.AsEnte()
	persona := titular.AsPersona()

	if ente != nil {
		return EnteFromDomain(ente)
	}

	if persona != nil {
		return PersonaFromDomain(persona)
	}

	return nil

}

func SujetoFromDomain(sujeto sujeto.Sujeto) Sujeto {
	if sujeto == nil {
		return nil
	}

	ente := sujeto.AsEnte()
	persona := sujeto.AsPersona()

	if ente != nil {
		return EnteFromDomain(ente)
	}

	if persona != nil {
		return PersonaFromDomain(persona)
	}

	return nil

}

func EnteFromDomain(ente *sujeto.Ente) *Ente {
	if ente == nil {
		return nil
	}

	return &Ente{
		ID:            ente.ID().String(),
		RazonSocial:   ente.RazonSocial(),
		Email:         ente.Email().String(),
		Telefono:      ente.Telefono().String(),
		Cedula:        ente.Cedula().String(),
		Registro:      ente.Audit().CreatedAt,
		Actualizacion: ente.Audit().UpdatedAt,
		Representante: &Persona{},
		DisplayName:   ente.DisplayName(),
	}

}

func PersonaFromDomain(persona *sujeto.Persona) *Persona {
	if persona == nil {
		return nil
	}

	return &Persona{
		ID:            persona.ID().String(),
		Nombres:       persona.Nombres(),
		Apellidos:     persona.Apellidos(),
		Email:         persona.Email().String(),
		Telefono:      persona.Telefono().String(),
		Cedula:        persona.Cedula().String(),
		Registro:      persona.Audit().CreatedAt,
		Actualizacion: persona.Audit().UpdatedAt,
		DisplayName:   persona.DisplayName(),
	}

}
