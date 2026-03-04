package command_test

import (
	"strings"
	"testing"
	"time"

	"github.com/Sanaruca/condominio/internal/administracion/app/command"
	"github.com/Sanaruca/condominio/internal/core/testing/utils"
	"github.com/Sanaruca/condominio/internal/pagos/types/moneda"
)

func TestValidateRegistrarGastoDTO(t *testing.T) {
	tests := []utils.ValidableTest{
		{
			Name: "Descripcion nula",
			Input: &command.RegistrarGastoDTO{
				Proveedor: "pvdr0",
				Monto:     1000,
				Moneda:    moneda.USD,
				Tasa:      36,
				// Descripcion: nil (omitted)
			},
			WantErr: false,
		},
		{
			Name: "Fecha nula",
			Input: &command.RegistrarGastoDTO{
				Proveedor: "pvdr0",
				Monto:     1000,
				Moneda:    moneda.USD,
				Tasa:      36,
				// Fecha: nil (omitted)
			},
			WantErr: false, // si fecha es opcional (usa hoy por default)
		},
		{
			Name: "Descripcion exactamente 500 chars",
			Input: &command.RegistrarGastoDTO{
				Proveedor:   "pvdr0",
				Monto:       1000,
				Moneda:      moneda.USD,
				Tasa:        36,
				Descripcion: utils.StringPtr(strings.Repeat("a", 500)),
			},
			WantErr: false,
		},
		{
			Name: "Descripcion 501 chars",
			Input: &command.RegistrarGastoDTO{
				Proveedor:   "pvdr0",
				Monto:       1000,
				Moneda:      moneda.USD,
				Tasa:        36,
				Descripcion: utils.StringPtr(strings.Repeat("a", 501)),
			},
			WantErr: true,
			ErrMsg:  "La descripción no puede exceder los 500 caracteres",
		},
		// TODO: Implementar validación de proveedor muy largo
		// {
		// 	name: "Proveedor muy largo",
		// 	Input: command.RegistrarGastoDTO{
		// 		Proveedor: strings.Repeat("a", 256),
		// 	},
		// 	WantErr: true,
		// 	ErrMsg:  "El proveedor no puede exceder X caracteres",
		// },
		// {
		// 	name: "Proveedor solo espacios",
		// 	Input: command.RegistrarGastoDTO{
		// 		Proveedor: "   ",
		// 		// resto válido
		// 	},
		// 	WantErr: true,
		// 	ErrMsg:  "El proveedor no puede estar vacío",
		// },
		{
			Name: "Moneda vacia",
			Input: &command.RegistrarGastoDTO{
				Proveedor: "pvdr0",
				Monto:     1000,
				Moneda:    "",
				Tasa:      36,
			},
			WantErr: true,
			ErrMsg:  "La moneda es requerida",
		},
		{
			Name: "Fecha exactamente hoy",
			Input: &command.RegistrarGastoDTO{
				Proveedor: "pvdr0",
				Monto:     1000,
				Moneda:    moneda.USD,
				Tasa:      36,
				Fecha:     utils.TimePtr(time.Now().Truncate(24 * time.Hour)),
			},
			WantErr: false,
		},
		{
			Name: "RegistrarGastoDTO válido",
			Input: &command.RegistrarGastoDTO{
				Proveedor:   "pvdr0",
				Monto:       1000,
				Moneda:      moneda.USD,
				Tasa:        36,
				Fecha:       utils.TimePtr(time.Now().Add(-24 * time.Hour)), // Ayer
				Descripcion: utils.StringPtr("Descripción válida"),
			},
			WantErr: false,
		},
		{
			Name: "Proveedor vacío",
			Input: &command.RegistrarGastoDTO{
				Proveedor:   "",
				Monto:       1000,
				Moneda:      moneda.USD,
				Tasa:        36,
				Descripcion: utils.StringPtr("Descripción"),
			},
			WantErr: true,
			ErrMsg:  "El proveedor es requerido",
		},
		{
			Name: "Monto invalido - cero (0)",
			Input: &command.RegistrarGastoDTO{
				Proveedor:   "pvdr0",
				Monto:       0,
				Moneda:      moneda.USD,
				Tasa:        36,
				Descripcion: utils.StringPtr("Descripción"),
			},
			WantErr: true,
			ErrMsg:  "El monto es requerido",
		},
		{
			Name: "Monto invalido - valor negativo",
			Input: &command.RegistrarGastoDTO{
				Proveedor:   "pvdr0",
				Monto:       -100,
				Moneda:      moneda.USD,
				Tasa:        36,
				Descripcion: utils.StringPtr("Descripción"),
			},
			WantErr: true,
			ErrMsg:  "El monto debe ser mayor a 0",
		},
		{
			Name: "Moneda invalida",
			Input: &command.RegistrarGastoDTO{
				Proveedor:   "pvdr0",
				Monto:       1000,
				Moneda:      "MXN",
				Tasa:        12,
				Descripcion: utils.StringPtr("Descripción"),
			},
			WantErr: true,
			ErrMsg:  "'MXN' no es una moneda valida",
		},
		{
			Name: "Tasa invalida - cero (0)",
			Input: &command.RegistrarGastoDTO{
				Proveedor:   "pvdr0",
				Monto:       1000,
				Moneda:      moneda.USD,
				Tasa:        0,
				Descripcion: utils.StringPtr("Descripción"),
			},
			WantErr: true,
			ErrMsg:  "La tasa es requerida",
		},
		{
			Name: "Tasa invalida - valor negativo (-1)",
			Input: &command.RegistrarGastoDTO{
				Proveedor:   "pvdr0",
				Monto:       1000,
				Moneda:      moneda.USD,
				Tasa:        -1,
				Descripcion: utils.StringPtr("Descripción"),
			},
			WantErr: true,
			ErrMsg:  "La tasa debe ser mayor a 0",
		},
		{
			Name: "Fecha futura",
			Input: &command.RegistrarGastoDTO{
				Proveedor:   "pvdr0",
				Monto:       1000,
				Moneda:      moneda.USD,
				Tasa:        36,
				Fecha:       utils.TimePtr(time.Now().Add(24 * time.Hour)), // Mañana
				Descripcion: utils.StringPtr("Descripción"),
			},
			WantErr: true,
			ErrMsg:  "La fecha no puede ser futura",
		},
		{
			Name: "Descripción vacía",
			Input: &command.RegistrarGastoDTO{
				Proveedor:   "pvdr0",
				Monto:       1000,
				Moneda:      moneda.USD,
				Tasa:        36,
				Descripcion: utils.StringPtr(""),
			},
			WantErr: true,
			ErrMsg:  "La descripción no puede estar vacía",
		},
		{
			Name: "Descripción muy larga",
			Input: &command.RegistrarGastoDTO{
				Proveedor: "pvdr0",
				Monto:     1000,
				Moneda:    moneda.USD,
				Tasa:      36,
				Descripcion: utils.StringPtr(
					`Esta es una descripción extremadamente larga que excede el límite de 500 caracteres.
				Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nullam auctor, nisl eget ultricies tincidunt, 
				nunc nisl aliquam nunc, vitae aliquam nisl nunc vitae nisl. Sed vitae nisl eget nisl aliquam tincidunt.
				Nullam auctor, nisl eget ultricies tincidunt, nunc nisl aliquam nunc, vitae aliquam nisl nunc vitae nisl.
				Sed vitae nisl eget nisl aliquam tincidunt. Nullam auctor, nisl eget ultricies tincidunt, nunc nisl
				aliquam nunc, vitae aliquam nisl nunc vitae nisl. Sed vitae nisl eget nisl aliquam tincidunt.`,
				),
			},
			WantErr: true,
			ErrMsg:  "La descripción no puede exceder los 500 caracteres",
		},
	}

	utils.TestValidables(t, tests, true)

}
