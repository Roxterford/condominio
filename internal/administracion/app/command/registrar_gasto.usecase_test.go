package command_test

import (
	"strings"
	"testing"
	"time"

	"github.com/Sanaruca/condominio/internal/administracion"
	"github.com/Sanaruca/condominio/internal/administracion/app/command"
	"github.com/Sanaruca/condominio/internal/administracion/app/query"
	"github.com/Sanaruca/condominio/internal/administracion/types/tipodeproveedor"
	"github.com/Sanaruca/condominio/internal/core/context"
	"github.com/Sanaruca/condominio/internal/core/session"
	"github.com/Sanaruca/condominio/internal/core/testing/utils"
	"github.com/Sanaruca/condominio/internal/pagos/types/moneda"
	"github.com/Sanaruca/condominio/internal/usuarios"
	"github.com/stretchr/testify/assert"
)

func TestRegistrarGasto(t *testing.T) {
	db := utils.SetupInMemoryTestDB(t)

	obtenerGasto := query.NewObtenerGasto()
	uc := command.NewRegistrarGasto(obtenerGasto)

	test_user := &usuarios.UsuarioPayload{
		ID:    "user0",
		Email: "test@example.com",
	}

	// Insertar proveedor de prueba
	if err := db.Create(&administracion.Proveedor{
		ID:            "pvdr0",
		Nombre:        "Proveedor de Prueba",
		Tipo:          tipodeproveedor.PersonaNatural,
		Email:         "proveedor@prueba.com",
		Telefono:      "123456789",
		Direccion:     stringPtr("Calle de Prueba 123"),
		Registro:      time.Now(),
		Actualizacion: time.Now(),
	}).Error; err != nil {
		t.Fatal("error insertando proveedor")
	}

	ctx := context.GetAdminContext(
		context.New(t.Context(), db).
			WithSession(session.New(test_user)),
	)

	dto := command.RegistrarGastoDTO{
		Proveedor: "pvdr0",
		Monto:     10_00,
		Moneda:    moneda.USD,
		Tasa:      10_00,
	}

	gasto, err := uc.Exec(ctx, dto)

	assert.NoError(t, err)
	assert.NotNil(t, gasto)

	if gasto.Proveedor != dto.Proveedor ||
		gasto.Monto != dto.Monto ||
		gasto.Moneda != string(dto.Moneda) ||
		gasto.Tasa != dto.Tasa ||
		gasto.Descripcion != dto.Descripcion {
		t.Errorf(
			"El input no coincide con el gasto registrado:\nInput:\n%s\nResultado:\n%s",
			utils.ToJSON(dto),
			utils.ToJSON(gasto),
		)
	}

	if gasto.RegistradoPor != test_user.ID {
		t.Errorf(
			"El usuario registrador no coincide con el usuario actual:\nUsuario registrador: %s\nUsuario actual: %s",
			gasto.RegistradoPor,
			test_user.ID,
		)
	}

	if gasto.ActualizadoPor != test_user.ID {
		t.Errorf(
			"El usuario actualizador no coincide con el usuario actual:\nUsuario actualizador: %s\nUsuario actual: %s",
			gasto.ActualizadoPor,
			test_user.ID,
		)
	}
}

func TestRegistrarGasto_Total(t *testing.T) {

	db := utils.SetupInMemoryTestDB(t)

	obtenerGasto := query.NewObtenerGasto()
	uc := command.NewRegistrarGasto(obtenerGasto)

	test_user := &usuarios.UsuarioPayload{
		ID:    "user0",
		Email: "test@example.com",
	}

	ctx := context.GetAdminContext(
		context.New(t.Context(), db).
			WithSession(session.New(test_user)),
	)

	// Insertar proveedor de prueba
	if err := db.Create(&administracion.Proveedor{
		ID:            "pvdr0",
		Nombre:        "Proveedor de Prueba",
		Tipo:          tipodeproveedor.PersonaNatural,
		Email:         "proveedor@prueba.com",
		Telefono:      "123456789",
		Direccion:     stringPtr("Calle de Prueba 123"),
		Registro:      time.Now(),
		Actualizacion: time.Now(),
	}).Error; err != nil {
		t.Fatal("error insertando proveedor")
	}

	tests := []struct {
		name     string
		monto    int
		moneda   moneda.Moneda
		tasa     int
		expected int
	}{
		// Tests en USD (directo)
		{
			name:     "Tasa 10 VED = 1 USD, Gasto 25 USD -> Total 25 USD",
			monto:    25_00,
			moneda:   moneda.USD,
			tasa:     10_00,
			expected: 25_00,
		},
		{
			name:     "Tasa 15 VED = 1 USD, Gasto 7.5 USD -> Total 7.5 USD",
			monto:    7_50,
			moneda:   moneda.USD,
			tasa:     15_00,
			expected: 7_50,
		},
		{
			name:     "Tasa 10.00 VED = 1 USD, Gasto 25.50 USD -> Total 25.50 USD",
			monto:    2550,
			moneda:   moneda.USD,
			tasa:     1000,
			expected: 2550,
		},
		{
			name:     "Tasa 15.75 VED = 1 USD, Gasto 7.25 USD -> Total 7.25 USD",
			monto:    725,
			moneda:   moneda.USD,
			tasa:     1575,
			expected: 725,
		},
		{
			name:     "Tasa 25.99 VED = 1 USD, Gasto 100.75 USD -> Total 100.75 USD",
			monto:    10075,
			moneda:   moneda.USD,
			tasa:     2599,
			expected: 10075,
		},
		// Tests en VED (conversión)
		{
			name:     "Tasa 10 VED = 1 USD, Gasto 100 VED -> Total 10 USD",
			monto:    100_00,
			moneda:   moneda.VED,
			tasa:     10_00,
			expected: 10_00,
		},
		{
			name:     "Tasa 5 VED = 1 USD, Gasto 25 VED -> Total 5 USD",
			monto:    25_00,
			moneda:   moneda.VED,
			tasa:     5_00,
			expected: 5_00,
		},
		{
			name:     "Tasa 1 VED = 1 USD, Gasto 10 VED -> Total 10 USD",
			monto:    10_00,
			moneda:   moneda.VED,
			tasa:     1_00,
			expected: 10_00,
		},
		{
			name:     "Tasa 20 VED = 1 USD, Gasto 40 VED -> Total 2 USD",
			monto:    40_00,
			moneda:   moneda.VED,
			tasa:     20_00,
			expected: 2_00,
		},
		{
			name:     "Tasa 10 VED = 1 USD, Gasto 50 VED -> Total 5 USD",
			monto:    50_00,
			moneda:   moneda.VED,
			tasa:     10_00,
			expected: 5_00,
		},
		{
			name:     "Tasa 2.5 VED = 1 USD, Gasto 12.5 VED -> Total 5 USD",
			monto:    12_50,
			moneda:   moneda.VED,
			tasa:     2_50,
			expected: 5_00,
		},
		{
			name:     "Tasa 10.50 VED = 1 USD, Gasto 105.00 VED -> Total 10.00 USD",
			monto:    105_00,
			moneda:   moneda.VED,
			tasa:     10_50,
			expected: 10_00, // (10500 * 100) / 1050 = 1000000 / 1050 = 1000 (10.00 USD)
		},
		{
			name:     "Tasa 7.25 VED = 1 USD, Gasto 36.25 VED -> Total 5.00 USD",
			monto:    36_25,
			moneda:   moneda.VED,
			tasa:     7_25,
			expected: 5_00, // (3625 * 100) / 725 = 362500 / 725 = 500 (5.00 USD)
		},
		{
			name:     "Tasa 12.50 VED = 1 USD, Gasto 25.00 VED -> Total 2.00 USD",
			monto:    25_00,
			moneda:   moneda.VED,
			tasa:     12_50,
			expected: 2_00, // (2500 * 100) / 1250 = 250000 / 1250 = 200 (2.00 USD)
		},
		{
			name:     "Tasa 33.33 VED = 1 USD, Gasto 99.99 VED -> Total 3.00 USD",
			monto:    99_99,
			moneda:   moneda.VED,
			tasa:     33_33,
			expected: 300, // (9999 * 100) / 3333 = 999900 / 3333 ≈ 300 (redondeo entero)
		},
		{
			name:     "Tasa 5.99 VED = 1 USD, Gasto 11.98 VED -> Total 2.00 USD",
			monto:    11_98,
			moneda:   moneda.VED,
			tasa:     5_99,
			expected: 2_00, // (1198 * 100) / 599 = 119800 / 599 = 200 (2.00 USD)
		},
	}

	for _, tt := range tests {
		test := tt

		t.Run(test.name, func(t *testing.T) {

			dto := command.RegistrarGastoDTO{
				Proveedor: "pvdr0",
				Monto:     test.monto,
				Moneda:    test.moneda,
				Tasa:      test.tasa,
			}

			gasto, err := uc.Exec(ctx, dto)

			if err != nil {
				t.Fatalf("error al ejecutar el caso de uso: %v", err)
			}
			assert.NotNil(t, gasto)
			assert.Equal(t, test.expected, gasto.Total, "Falló la prueba: %s", test.name)
		})

	}

}

func TestRegistrarGasto_ProveedorNoEncontrado(t *testing.T) {
	db := utils.SetupInMemoryTestDB(t)

	obtenerGasto := query.NewObtenerGasto()
	uc := command.NewRegistrarGasto(obtenerGasto)

	test_user := &usuarios.UsuarioPayload{
		ID:    "user0",
		Email: "test@example.com",
	}

	ctx := context.GetAdminContext(
		context.New(t.Context(), db).
			WithSession(session.New(test_user)),
	)

	dto := command.RegistrarGastoDTO{
		Proveedor: "pvdr0",
		Monto:     10_00,
		Moneda:    moneda.USD,
		Tasa:      10_00,
	}

	_, err := uc.Exec(ctx, dto)

	assert.Error(t, err)
	assert.Equal(t, administracion.ErrProveedorNoEncontrado, err)
}

func TestValidateRegistrarGastoDTO(t *testing.T) {
	tests := []struct {
		name    string
		dto     command.RegistrarGastoDTO
		wantErr bool
		errMsg  string
	}{
		{
			name: "Descripcion nula",
			dto: command.RegistrarGastoDTO{
				Proveedor: "pvdr0",
				Monto:     1000,
				Moneda:    moneda.USD,
				Tasa:      36,
				// Descripcion: nil (omitted)
			},
			wantErr: false,
		},
		{
			name: "Fecha nula",
			dto: command.RegistrarGastoDTO{
				Proveedor: "pvdr0",
				Monto:     1000,
				Moneda:    moneda.USD,
				Tasa:      36,
				// Fecha: nil (omitted)
			},
			wantErr: false, // si fecha es opcional (usa hoy por default)
		},
		{
			name: "Descripcion exactamente 500 chars",
			dto: command.RegistrarGastoDTO{
				Proveedor:   "pvdr0",
				Monto:       1000,
				Moneda:      moneda.USD,
				Tasa:        36,
				Descripcion: stringPtr(strings.Repeat("a", 500)),
			},
			wantErr: false,
		},
		{
			name: "Descripcion 501 chars",
			dto: command.RegistrarGastoDTO{
				Proveedor:   "pvdr0",
				Monto:       1000,
				Moneda:      moneda.USD,
				Tasa:        36,
				Descripcion: stringPtr(strings.Repeat("a", 501)),
			},
			wantErr: true,
			errMsg:  "La descripción no puede exceder los 500 caracteres",
		},
		// TODO: Implementar validación de proveedor muy largo
		// {
		// 	name: "Proveedor muy largo",
		// 	dto: command.RegistrarGastoDTO{
		// 		Proveedor: strings.Repeat("a", 256),
		// 	},
		// 	wantErr: true,
		// 	errMsg:  "El proveedor no puede exceder X caracteres",
		// },
		// {
		// 	name: "Proveedor solo espacios",
		// 	dto: command.RegistrarGastoDTO{
		// 		Proveedor: "   ",
		// 		// resto válido
		// 	},
		// 	wantErr: true,
		// 	errMsg:  "El proveedor no puede estar vacío",
		// },
		{
			name: "Moneda vacia",
			dto: command.RegistrarGastoDTO{
				Proveedor: "pvdr0",
				Monto:     1000,
				Moneda:    "",
				Tasa:      36,
			},
			wantErr: true,
			errMsg:  "La moneda es requerida",
		},
		{
			name: "Fecha exactamente hoy",
			dto: command.RegistrarGastoDTO{
				Proveedor: "pvdr0",
				Monto:     1000,
				Moneda:    moneda.USD,
				Tasa:      36,
				Fecha:     timePtr(time.Now().Truncate(24 * time.Hour)),
			},
			wantErr: false,
		},
		{
			name: "RegistrarGastoDTO válido",
			dto: command.RegistrarGastoDTO{
				Proveedor:   "pvdr0",
				Monto:       1000,
				Moneda:      moneda.USD,
				Tasa:        36,
				Fecha:       timePtr(time.Now().Add(-24 * time.Hour)), // Ayer
				Descripcion: stringPtr("Descripción válida"),
			},
			wantErr: false,
		},
		{
			name: "Proveedor vacío",
			dto: command.RegistrarGastoDTO{
				Proveedor:   "",
				Monto:       1000,
				Moneda:      moneda.USD,
				Tasa:        36,
				Descripcion: stringPtr("Descripción"),
			},
			wantErr: true,
			errMsg:  "El proveedor es requerido",
		},
		{
			name: "Monto invalido - cero (0)",
			dto: command.RegistrarGastoDTO{
				Proveedor:   "pvdr0",
				Monto:       0,
				Moneda:      moneda.USD,
				Tasa:        36,
				Descripcion: stringPtr("Descripción"),
			},
			wantErr: true,
			errMsg:  "El monto es requerido",
		},
		{
			name: "Monto invalido - valor negativo",
			dto: command.RegistrarGastoDTO{
				Proveedor:   "pvdr0",
				Monto:       -100,
				Moneda:      moneda.USD,
				Tasa:        36,
				Descripcion: stringPtr("Descripción"),
			},
			wantErr: true,
			errMsg:  "El monto debe ser mayor a 0",
		},
		{
			name: "Moneda invalida",
			dto: command.RegistrarGastoDTO{
				Proveedor:   "pvdr0",
				Monto:       1000,
				Moneda:      "MXN",
				Tasa:        12,
				Descripcion: stringPtr("Descripción"),
			},
			wantErr: true,
			errMsg:  "'MXN' no es una moneda valida",
		},
		{
			name: "Tasa invalida - cero (0)",
			dto: command.RegistrarGastoDTO{
				Proveedor:   "pvdr0",
				Monto:       1000,
				Moneda:      moneda.USD,
				Tasa:        0,
				Descripcion: stringPtr("Descripción"),
			},
			wantErr: true,
			errMsg:  "La tasa es requerida",
		},
		{
			name: "Tasa invalida - valor negativo (-1)",
			dto: command.RegistrarGastoDTO{
				Proveedor:   "pvdr0",
				Monto:       1000,
				Moneda:      moneda.USD,
				Tasa:        -1,
				Descripcion: stringPtr("Descripción"),
			},
			wantErr: true,
			errMsg:  "La tasa debe ser mayor a 0",
		},
		{
			name: "Fecha futura",
			dto: command.RegistrarGastoDTO{
				Proveedor:   "pvdr0",
				Monto:       1000,
				Moneda:      moneda.USD,
				Tasa:        36,
				Fecha:       timePtr(time.Now().Add(24 * time.Hour)), // Mañana
				Descripcion: stringPtr("Descripción"),
			},
			wantErr: true,
			errMsg:  "La fecha no puede ser futura",
		},
		{
			name: "Descripción vacía",
			dto: command.RegistrarGastoDTO{
				Proveedor:   "pvdr0",
				Monto:       1000,
				Moneda:      moneda.USD,
				Tasa:        36,
				Descripcion: stringPtr(""),
			},
			wantErr: true,
			errMsg:  "La descripción no puede estar vacía",
		},
		{
			name: "Descripción muy larga",
			dto: command.RegistrarGastoDTO{
				Proveedor: "pvdr0",
				Monto:     1000,
				Moneda:    moneda.USD,
				Tasa:      36,
				Descripcion: stringPtr(
					`Esta es una descripción extremadamente larga que excede el límite de 500 caracteres.
				Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nullam auctor, nisl eget ultricies tincidunt, 
				nunc nisl aliquam nunc, vitae aliquam nisl nunc vitae nisl. Sed vitae nisl eget nisl aliquam tincidunt.
				Nullam auctor, nisl eget ultricies tincidunt, nunc nisl aliquam nunc, vitae aliquam nisl nunc vitae nisl.
				Sed vitae nisl eget nisl aliquam tincidunt. Nullam auctor, nisl eget ultricies tincidunt, nunc nisl
				aliquam nunc, vitae aliquam nisl nunc vitae nisl. Sed vitae nisl eget nisl aliquam tincidunt.`,
				),
			},
			wantErr: true,
			errMsg:  "La descripción no puede exceder los 500 caracteres",
		},
	}

	for _, tezt := range tests {
		test := tezt
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			err := test.dto.Validate()

			if err != nil {
				t.Log(err.Error())
			}

			if test.wantErr {
				assert.Error(t, err, "El test '%s' debería haber fallado", test.name)
				if test.errMsg != "" {
					assert.Contains(
						t,
						err.Error(),
						test.errMsg,
						"Se esperaba que el error contenga el mensaje: %s\nPero en su lugar fue: %s",
						test.errMsg,
						err.Error(),
					)
				} else {
					t.Log("Cambie el mensaje de error esperado en el test '" + test.name + "'")
					t.Fatal("Proporcione un mensaje de error al esperar uno")
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// Helper functions
func stringPtr(s string) *string {
	return &s
}

func timePtr(t time.Time) *time.Time {
	return &t
}
