package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPhoneFactory_New(t *testing.T) {
	factory := NewPhoneFactory([]string{"58"}, []string{"414", "424", "416", "412"})

	tests := []struct {
		name    string
		input   string
		wantErr bool
		checkFn func(*testing.T, Phone)
	}{
		{
			name:    "teléfono válido",
			input:   "58-414-1234567",
			wantErr: false,
			checkFn: func(t *testing.T, p Phone) {
				assert.Equal(t, "58", p.CountryCode())
				assert.Equal(t, "414", p.Provider())
				assert.Equal(t, "1234567", p.Subscriber())
			},
		},
		{
			name:    "teléfono con otro proveedor válido",
			input:   "58-416-9876543",
			wantErr: false,
			checkFn: func(t *testing.T, p Phone) {
				assert.Equal(t, "416", p.Provider())
			},
		},
		{
			name:    "teléfono inválido sin guiones",
			input:   "584141234567",
			wantErr: true,
		},
		{
			name:    "teléfono inválido con muchos segmentos",
			input:   "58-414-12-34567",
			wantErr: true,
		},
		{
			name:    "teléfono con muy pocos dígitos en suscriptor",
			input:   "58-414-12345",
			wantErr: true,
		},
		{
			name:    "país no permitido",
			input:   "57-414-1234567",
			wantErr: true,
		},
		{
			name:    "proveedor no permitido",
			input:   "58-411-1234567",
			wantErr: true,
		},
		{
			name:    "teléfono válido con espacios como separadores",
			input:   "58-414-1234567",
			wantErr: false,
			checkFn: func(t *testing.T, p Phone) {
				assert.Equal(t, "1234567", p.Subscriber())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, err := factory.New(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			tt.checkFn(t, p)
		})
	}
}

func TestPhone_FullNumber(t *testing.T) {
	factory := NewPhoneFactory([]string{"58"}, []string{"414"})

	p, err := factory.New("58-414-1234567")
	require.NoError(t, err)

	assert.Equal(t, "+584141234567", p.FullNumber())
}

func TestPhone_String(t *testing.T) {
	factory := NewPhoneFactory([]string{"58"}, []string{"414"})

	p, err := factory.New("58-414-1234567")
	require.NoError(t, err)

	assert.Equal(t, p.FullNumber(), p.String())
}

func TestPhoneFactory_WithoutRestrictions(t *testing.T) {
	factory := NewPhoneFactory([]string{}, []string{})

	p, err := factory.New("1-234-567890")
	require.NoError(t, err)
	assert.Equal(t, "1", p.CountryCode())
	assert.Equal(t, "234", p.Provider())
	assert.Equal(t, "567890", p.Subscriber())
}
