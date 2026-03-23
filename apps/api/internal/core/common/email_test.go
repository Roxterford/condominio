package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEmailFactory_New(t *testing.T) {
	factory := NewEmailFactory([]string{"gmail.com", "example.com", "test.co", "sub.example.com"})

	tests := []struct {
		name    string
		input   string
		wantErr bool
		checkFn func(*testing.T, Email)
	}{
		{
			name:    "email válido",
			input:   "test@example.com",
			wantErr: false,
			checkFn: func(t *testing.T, e Email) {
				assert.Equal(t, "test@example.com", e.Address())
			},
		},
		{
			name:    "email válido con mayúsculas",
			input:   "Test@Example.com",
			wantErr: false,
			checkFn: func(t *testing.T, e Email) {
				assert.Equal(t, "test@example.com", e.Address())
			},
		},
		{
			name:    "email válido con punto en dominio",
			input:   "test@sub.example.com",
			wantErr: false,
			checkFn: func(t *testing.T, e Email) {
				assert.Equal(t, "test@sub.example.com", e.Address())
			},
		},
		{
			name:    "email con dominio permitido",
			input:   "user@gmail.com",
			wantErr: false,
			checkFn: func(t *testing.T, e Email) {
				assert.Equal(t, "user@gmail.com", e.Address())
			},
		},
		{
			name:    "email inválido sin @",
			input:   "testexample.com",
			wantErr: true,
		},
		{
			name:    "email inválido sin dominio",
			input:   "test@",
			wantErr: true,
		},
		{
			name:    "email inválido sin usuario",
			input:   "@example.com",
			wantErr: true,
		},
		{
			name:    "email inválido con espacio",
			input:   "test @example.com",
			wantErr: true,
		},
		{
			name:    "dominio no permitido",
			input:   "test@prohibited.com",
			wantErr: true,
		},
		{
			name:    "email con dominio permitido test.co",
			input:   "user@test.co",
			wantErr: false,
			checkFn: func(t *testing.T, e Email) {
				assert.Equal(t, "user@test.co", e.Address())
			},
		},
		{
			name:    "email con múltiples puntos",
			input:   "test.mail@example.com",
			wantErr: false,
			checkFn: func(t *testing.T, e Email) {
				assert.Equal(t, "test.mail@example.com", e.Address())
			},
		},
		{
			name:    "email con guión",
			input:   "test-user@example.com",
			wantErr: false,
			checkFn: func(t *testing.T, e Email) {
				assert.Equal(t, "test-user@example.com", e.Address())
			},
		},
		{
			name:    "email con guión bajo",
			input:   "test_user@example.com",
			wantErr: false,
			checkFn: func(t *testing.T, e Email) {
				assert.Equal(t, "test_user@example.com", e.Address())
			},
		},
		{
			name:    "email con más",
			input:   "test+mail@example.com",
			wantErr: false,
			checkFn: func(t *testing.T, e Email) {
				assert.Equal(t, "test+mail@example.com", e.Address())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e, err := factory.New(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			tt.checkFn(t, e)
		})
	}
}

func TestEmail_Equals(t *testing.T) {
	factory := NewEmailFactory([]string{})

	email1, _ := factory.New("test@example.com")
	email2, _ := factory.New("TEST@EXAMPLE.COM")
	email3, _ := factory.New("other@example.com")

	tests := []struct {
		name     string
		email1   Email
		email2   Email
		expected bool
	}{
		{"mismo email, diferente caso", email1, email2, true},
		{"diferentes emails", email1, email3, false},
		{"mismo email", email1, email1, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.email1.Equals(tt.email2))
		})
	}
}

func TestEmail_String(t *testing.T) {
	factory := NewEmailFactory([]string{})

	email, err := factory.New("test@example.com")
	require.NoError(t, err)

	assert.Equal(t, "test@example.com", email.String())
}

func TestEmailFactory_WithoutRestrictions(t *testing.T) {
	factory := NewEmailFactory([]string{})

	email, err := factory.New("any@domain.com")
	require.NoError(t, err)
	assert.Equal(t, "any@domain.com", email.Address())
}

func TestEmailFactory_CaseInsensitiveDomain(t *testing.T) {
	factory := NewEmailFactory([]string{"GMAIL.COM"})

	email, err := factory.New("test@Gmail.com")
	require.NoError(t, err)
	assert.Equal(t, "test@gmail.com", email.Address())
}
