package sujeto_test

import (
	"testing"

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
	)

	got := p.DisplayName()
	if got != "Ana" {
		t.Fatalf("DisplayName() = %q, esperado %q", got, "Ana")
	}
}
