package sujeto_test

import (
	"testing"
	"time"

	"github.com/Sanaruca/condominio/internal/core/common"
	"github.com/Sanaruca/condominio/internal/unidades/models/sujeto"
)

func nuevoFactory() *sujeto.SujetoFactory {
	return sujeto.NewSujetoFactory(
		common.NewEmailFactory(nil),
		common.NewPhoneFactory(nil, nil),
	)
}

func TestPersonaDisplayNamePrimerNombreYPrimerApellido(t *testing.T) {
	p := nuevoFactory().AssemblePersona(
		"1",
		"V-12345678",
		"Juan Carlos",
		"Pérez Gómez",
		"juan@example.com",
		"04120000000",
		time.Now(),
	)

	got := p.DisplayName()
	if got != "Juan Pérez" {
		t.Fatalf("DisplayName() = %q, esperado %q", got, "Juan Pérez")
	}
}

func TestPersonaDisplayNameSoloNombres(t *testing.T) {
	p := nuevoFactory().AssemblePersona(
		"1",
		"V-12345678",
		"Ana María",
		"",
		"ana@example.com",
		"04120000000",
		time.Now(),
	)

	got := p.DisplayName()
	if got != "Ana" {
		t.Fatalf("DisplayName() = %q, esperado %q", got, "Ana")
	}
}

func TestPersonaAssembleRehidrataRegistro(t *testing.T) {
	registro := time.Date(2024, 5, 1, 10, 30, 0, 0, time.UTC)

	p := nuevoFactory().AssemblePersona(
		"1",
		"V-12345678",
		"Juan",
		"Pérez",
		"juan@example.com",
		"04120000000",
		registro,
	)

	got := p.Audit().CreatedAt
	if !got.Equal(registro) {
		t.Fatalf("Audit().CreatedAt = %v, esperado %v", got, registro)
	}
}
