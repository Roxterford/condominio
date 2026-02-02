package command_test

import (
	"testing"

	"github.com/Sanaruca/condominio/internal/administracion/app/command"
	"github.com/Sanaruca/condominio/internal/administracion/types/tipodeproveedor"
	"github.com/Sanaruca/condominio/internal/core/testing/utils"
)

// TODO: Refactorizar partiendo de un input valido e ir cambiando los campos
// para probar cada validacion, reduciendo asi las lineas de codigo. Se puede
// crear un RegistrarProveedorDTO valido e ir creando copias con los campos
// modificados segun cada caso de prueba.
// Ejemplo:
//
//	validDTO := &command.RegistrarProveedorDTO{
//		Rif:      "J-123456789",
//		Nombre:   "Proveedor S.A.",
//		Tipo:     tipodeproveedor.Compania,
//		Email:    "proveedor@ejemplo.com",
//		Telefono: "04121234567",
//	}
//
//	nombre_vacio := validDTO
//	nombre_vacio.Nombre = ""
func TestValidateRegistrarProveedorDTO(t *testing.T) {

	tests := []utils.ValidableTest{
		// Casos validos
		{
			Name: "RegistrarProveedorDTO valido",
			Input: &command.RegistrarProveedorDTO{
				Rif:      "J-123456789",
				Nombre:   "Proveedor S.A.",
				Tipo:     tipodeproveedor.Compania,
				Email:    "proveedor@ejemplo.com",
				Telefono: "04121234567",
			},
			WantErr: false,
		},
		// Casos con datos requeridos
		{
			Name: "RegistrarProveedorDTO con rif vacio",
			Input: &command.RegistrarProveedorDTO{
				Rif:      "", // <--
				Nombre:   "Proveedor S.A.",
				Tipo:     tipodeproveedor.Compania,
				Email:    "proveedor@ejemplo.com",
				Telefono: "04121234567",
			},
			WantErr: true,
			ErrMsg:  "El RIF es requerido",
		},
		{
			Name: "RegistrarProveedorDTO con nombre vacio",
			Input: &command.RegistrarProveedorDTO{
				Rif:      "J-123456789",
				Nombre:   "", // <--
				Tipo:     tipodeproveedor.Compania,
				Email:    "proveedor@ejemplo.com",
				Telefono: "04121234567",
			},
			WantErr: true,
			ErrMsg:  "El nombre es requerido",
		},
		{
			Name: "RegistrarProveedorDTO con tipo vacio",
			Input: &command.RegistrarProveedorDTO{
				Rif:      "J-123456789",
				Nombre:   "Proveedor S.A.",
				Tipo:     "", // <--
				Email:    "proveedor@ejemplo.com",
				Telefono: "04121234567",
			},
			WantErr: true,
			ErrMsg:  "El tipo es requerido",
		},
		{
			Name: "RegistrarProveedorDTO con email vacio",
			Input: &command.RegistrarProveedorDTO{
				Rif:      "J-123456789",
				Nombre:   "Proveedor S.A.",
				Tipo:     tipodeproveedor.Compania,
				Email:    "", // <--
				Telefono: "04121234567",
			},
			WantErr: true,
			ErrMsg:  "El email es requerido",
		},
		{
			Name: "RegistrarProveedorDTO con telefono vacio",
			Input: &command.RegistrarProveedorDTO{
				Rif:      "J-123456789",
				Nombre:   "Proveedor S.A.",
				Tipo:     tipodeproveedor.Compania,
				Email:    "proveedor@ejemplo.com",
				Telefono: "", // <--
			},
			WantErr: true,
			ErrMsg:  "El telefono es requerido",
		},
		// Email invalido
		{
			Name: "RegistrarProveedorDTO con email invalido",
			Input: &command.RegistrarProveedorDTO{
				Rif:      "J-123456789",
				Nombre:   "Proveedor S.A.",
				Tipo:     tipodeproveedor.Compania,
				Email:    "proveedor@ejemplo", // <--
				Telefono: "04121234567",
			},
			WantErr: true,
			ErrMsg:  "El email es invalido",
		},
		// Tipo de proveedor invalido
		{
			Name: "RegistrarProveedorDTO con tipo de proveedor invalido",
			Input: &command.RegistrarProveedorDTO{
				Rif:      "J-123456789",
				Nombre:   "Proveedor S.A.",
				Tipo:     "INVALIDO", // <--
				Email:    "proveedor@ejemplo.com",
				Telefono: "04121234567",
			},
			WantErr: true,
			ErrMsg:  "'INVALIDO' no es un tipo de proveedor valido",
		},
		// Campo Direccion opcional
		{
			Name: "RegistrarProveedorDTO con direccion opcional",
			Input: &command.RegistrarProveedorDTO{
				Rif:       "J-123456789",
				Nombre:    "Proveedor S.A.",
				Tipo:      tipodeproveedor.Compania,
				Email:     "proveedor@ejemplo.com",
				Telefono:  "04121234567",
				Direccion: utils.StringPtr("123123"), // <--
			},
			WantErr: false,
		},
		{
			Name: "RegistrarProveedorDTO con direccion opcional",
			Input: &command.RegistrarProveedorDTO{
				Rif:       "J-123456789",
				Nombre:    "Proveedor S.A.",
				Tipo:      tipodeproveedor.Compania,
				Email:     "proveedor@ejemplo.com",
				Telefono:  "04121234567",
				Direccion: utils.StringPtr("     "), // <--
			},
			WantErr: true,
			ErrMsg:  "La direccion no puede estar vacia",
		},
		// Telefono con formato invalido
		{
			Name: "RegistrarProveedorDTO con telefono con formato invalido",
			Input: &command.RegistrarProveedorDTO{
				Rif:      "J-123456789",
				Nombre:   "Proveedor S.A.",
				Tipo:     tipodeproveedor.Compania,
				Email:    "proveedor@ejemplo.com",
				Telefono: "012312311223", // <--
			},
			WantErr: true,
			ErrMsg:  "El teléfono debe tener formato 04141234567",
		},
		{
			Name: "RegistrarProveedorDTO con telefono con formato invalido",
			Input: &command.RegistrarProveedorDTO{
				Rif:      "J-123456789",
				Nombre:   "Proveedor S.A.",
				Tipo:     tipodeproveedor.Compania,
				Email:    "proveedor@ejemplo.com",
				Telefono: "04231234567", // <--
			},
			WantErr: true,
			ErrMsg:  "El teléfono debe tener formato 04141234567",
		},
	}

	utils.TestValidables(t, tests, true)

}
