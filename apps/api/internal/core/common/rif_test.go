package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewRif(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
		checkFn func(*testing.T, Rif)
	}{
		{
			name:    "RIF válido tipo V",
			input:   "V12345678",
			wantErr: false,
			checkFn: func(t *testing.T, r Rif) {
				assert.Equal(t, "V12345678", r.String())
			},
		},
		{
			name:    "RIF válido tipo J",
			input:   "J123456789",
			wantErr: false,
			checkFn: func(t *testing.T, r Rif) {
				assert.Equal(t, "J123456789", r.String())
			},
		},
		{
			name:    "RIF válido tipo E",
			input:   "E12345678",
			wantErr: false,
			checkFn: func(t *testing.T, r Rif) {
				assert.Equal(t, "E12345678", r.String())
			},
		},
		{
			name:    "RIF válido tipo G",
			input:   "G12345678",
			wantErr: false,
			checkFn: func(t *testing.T, r Rif) {
				assert.Equal(t, "G12345678", r.String())
			},
		},
		{
			name:    "RIF válido tipo P",
			input:   "P12345678",
			wantErr: false,
			checkFn: func(t *testing.T, r Rif) {
				assert.Equal(t, "P12345678", r.String())
			},
		},
		{
			name:    "RIF válido con guiones",
			input:   "J-12345678-9",
			wantErr: false,
			checkFn: func(t *testing.T, r Rif) {
				assert.Equal(t, "J123456789", r.String())
			},
		},
		{
			name:    "RIF válido con guiones y minúsculas",
			input:   "v-12345678",
			wantErr: false,
			checkFn: func(t *testing.T, r Rif) {
				assert.Equal(t, "V12345678", r.String())
			},
		},
		{
			name:    "RIF inválido con prefijo incorrecto",
			input:   "X12345678",
			wantErr: true,
		},
		{
			name:    "RIF inválido muy corto",
			input:   "V123",
			wantErr: true,
		},
		{
			name:    "RIF inválido vacío",
			input:   "",
			wantErr: true,
		},
		{
			name:    "RIF inválido con solo prefijo",
			input:   "V",
			wantErr: true,
		},
		{
			name:    "RIF inválido con 7 dígitos",
			input:   "V1234567",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := NewRif(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			tt.checkFn(t, r)
		})
	}
}

func TestRif_EsPersonaNatural(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"V prefix es persona natural", "V12345678", true},
		{"E prefix es persona natural", "E12345678", true},
		{"J prefix no es persona natural", "J12345678", false},
		{"G prefix no es persona natural", "G12345678", false},
		{"P prefix no es persona natural", "P12345678", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := NewRif(tt.input)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, r.EsPersonaNatural())
		})
	}
}

func TestRif_EsEmpresa(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"J prefix es empresa", "J12345678", true},
		{"V prefix no es empresa", "V12345678", false},
		{"E prefix no es empresa", "E12345678", false},
		{"G prefix no es empresa", "G12345678", false},
		{"P prefix no es empresa", "P12345678", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := NewRif(tt.input)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, r.EsEmpresa())
		})
	}
}
