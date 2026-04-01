package command_test

import (
	"testing"

	"github.com/Sanaruca/condominio/internal/administracion/app/command"
	"github.com/stretchr/testify/assert"
)

func TestValidateRegistrarProveedorDTO(t *testing.T) {
	direccion := "Calle Falsa 123"
	emailValido := "test@ejemplo.com"

	tests := []struct {
		name    string
		dto     command.RegistrarProveedorDTO
		wantErr bool
	}{
		{
			name: "válido con todos los campos",
			dto: command.RegistrarProveedorDTO{
				Rif:       "J-12345678-9",
				Nombre:    "Proveedor Test",
				Email:     emailValido,
				Telefono:  "04141234567",
				Direccion: &direccion,
			},
			wantErr: false,
		},
		{
			name: "RIF inválido (sin prefijo válido)",
			dto: command.RegistrarProveedorDTO{
				Rif:       "X-12345678-9",
				Nombre:    "Proveedor Test",
				Email:     emailValido,
				Telefono:  "04141234567",
				Direccion: &direccion,
			},
			wantErr: true,
		},
		{
			name: "RIF muy corto",
			dto: command.RegistrarProveedorDTO{
				Rif:       "J",
				Nombre:    "Proveedor Test",
				Email:     emailValido,
				Telefono:  "04141234567",
				Direccion: &direccion,
			},
			wantErr: true,
		},
		{
			name: "Nombre vacío",
			dto: command.RegistrarProveedorDTO{
				Rif:       "J-12345678-9",
				Nombre:    "",
				Email:     emailValido,
				Telefono:  "04141234567",
				Direccion: &direccion,
			},
			wantErr: true,
		},
		{
			name: "Nombre muy largo",
			dto: command.RegistrarProveedorDTO{
				Rif:       "J-12345678-9",
				Nombre:    string(make([]byte, 101)),
				Email:     emailValido,
				Telefono:  "04141234567",
				Direccion: &direccion,
			},
			wantErr: true,
		},
		{
			name: "Email inválido",
			dto: command.RegistrarProveedorDTO{
				Rif:       "J-12345678-9",
				Nombre:    "Proveedor Test",
				Email:     "email-invalido",
				Telefono:  "04141234567",
				Direccion: &direccion,
			},
			wantErr: true,
		},
		{
			name: "Teléfono vacío",
			dto: command.RegistrarProveedorDTO{
				Rif:       "J-12345678-9",
				Nombre:    "Proveedor Test",
				Email:     emailValido,
				Telefono:  "",
				Direccion: &direccion,
			},
			wantErr: true,
		},
		{
			name: "Teléfono muy largo",
			dto: command.RegistrarProveedorDTO{
				Rif:       "J-12345678-9",
				Nombre:    "Proveedor Test",
				Email:     emailValido,
				Telefono:  string(make([]byte, 21)),
				Direccion: &direccion,
			},
			wantErr: true,
		},
		{
			name: "RIF vacío",
			dto: command.RegistrarProveedorDTO{
				Rif:       "",
				Nombre:    "Proveedor Test",
				Email:     emailValido,
				Telefono:  "04141234567",
				Direccion: &direccion,
			},
			wantErr: true,
		},
	}

	for _, tst := range tests {
		test := tst
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			err := test.dto.Validate()
			if test.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
