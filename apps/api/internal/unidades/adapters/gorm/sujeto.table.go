package gorm

import (
	"time"

	"github.com/Sanaruca/condominio/internal/unidades/models/sujeto"
)

type TipoDeSujeto string

const (
	PERSONA_NATURAL TipoDeSujeto = "PERSONA_NATURAL"
	ENTE_JURIDICO   TipoDeSujeto = "ENTE_JURIDICO"
)

type Sujeto struct {
	ID                 string `gorm:"primaryKey"`
	Tipo               TipoDeSujeto
	DocumentoIdentidad string
	Nombres            *string /// Persona
	Apellidos          *string /// Persona
	RazonSocial        *string /// Jurídico
	RepresentanteID    *string `gorm:"column:representante"`
	Email              string
	Telefono           string
	Registro           time.Time

	Representante *Sujeto `gorm:"foreignKey:RepresentanteID"`
}

func (s *Sujeto) TableName() string {
	return "sujetos"
}

func (s *Sujeto) ToDoaminPersona(factory *sujeto.SujetoFactory) sujeto.Persona {

	if s == nil {
		return sujeto.Persona{}
	}

	var nombres string
	var apellidos string

	if s.Nombres != nil {
		nombres = *s.Nombres
	}

	if s.Apellidos != nil {
		apellidos = *s.Apellidos
	}

	persona := factory.AssemblePersona(
		s.ID,
		s.DocumentoIdentidad,
		nombres,
		apellidos,
		s.Email,
		s.Telefono,
	)

	return *persona
}
