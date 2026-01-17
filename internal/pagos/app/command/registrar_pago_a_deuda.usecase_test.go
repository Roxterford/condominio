package command_test

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Sanaruca/condominio/internal/administracion"
	"github.com/Sanaruca/condominio/internal/core/context"
	testingutils "github.com/Sanaruca/condominio/internal/core/testing/utils"
	"github.com/Sanaruca/condominio/internal/core/utils"

	"github.com/Sanaruca/condominio/internal/pagos"
	"github.com/Sanaruca/condominio/internal/pagos/app/command"
	"github.com/Sanaruca/condominio/internal/pagos/types/metododepago"
	"github.com/Sanaruca/condominio/internal/pagos/types/moneda"
	"github.com/Sanaruca/condominio/internal/villas"
	"github.com/Sanaruca/condominio/internal/villas/types/estadodeuda"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegistrarPagoADeuda_PagoNoEncontrado(t *testing.T) {
	db, mock := testingutils.SetupMockTestDB(t)

	ctx := context.New(t.Context(), db)

	uc := command.NewRegistarPagoADeuda()

	dto := command.RegistarPagoADeudaDTO{"null"}

	mock.ExpectQuery(
		fmt.Sprintf(
			"SELECT (.+) FROM %s WHERE id = \\?",
			utils.GetMockTableRegex("pagos"),
		),
	).
		WithArgs(dto.PagoID).
		WillReturnRows(sqlmock.NewRows([]string{}))

	_, err := uc.Exec(ctx, dto)

	assert.Error(t, err)

	assert.True(
		t,
		errors.Is(err, pagos.ErrPagoNoEncontrado),
		"Se esperaba ErrPagoNoEncontrado, pero se obtuvo: %v",
		err,
	)
}

func TestRegistrarPagoADeuda_DeudaNoEncontrada(t *testing.T) {
	db, mock := testingutils.SetupMockTestDB(t)
	ctx := context.New(t.Context(), db)
	uc := command.NewRegistarPagoADeuda()

	dto := command.RegistarPagoADeudaDTO{PagoID: "p0"}

	mock.ExpectQuery(
		fmt.Sprintf(
			"SELECT (.+) FROM %s WHERE id = \\?",
			utils.GetMockTableRegex("pagos"),
		),
	).
		WithArgs(dto.PagoID).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "villa", "cuenta"}).
				AddRow(dto.PagoID, 1, 100),
		)

	mock.ExpectQuery(
		fmt.Sprintf(
			"SELECT (.+) FROM %s WHERE villa = \\? AND estado <> ?",
			utils.GetMockTableRegex("deudas"),
		),
	).
		WithArgs(1, estadodeuda.Pagada).
		WillReturnRows(sqlmock.NewRows([]string{}))

	res, err := uc.Exec(ctx, dto)

	assert.NoError(t, err)
	assert.Nil(t, res)
}

func TestRegistrarPagoADeuda(t *testing.T) {

	type expected struct {
		// destinado a deuda
		destinado int
		// saldo sobrante del pago
		cuenta int
	}

	tests := []struct {
		villa int
		pago  int
		deuda int

		expected expected
	}{
		// Caso 1: Pago cubre exactamente la deuda (100 USD)
		{villa: 1, pago: 100_00, deuda: 100_00, expected: expected{destinado: 100_00, cuenta: 0}},
		// Caso 2: Pago mayor que la deuda (200 USD vs 150 USD)
		{
			villa:    2,
			pago:     200_00,
			deuda:    150_00,
			expected: expected{destinado: 150_00, cuenta: 50_00},
		},
		// Caso 3: Pago menor que la deuda (80 USD vs 120 USD)
		{villa: 3, pago: 80_00, deuda: 120_00, expected: expected{destinado: 80_00, cuenta: 0}},
		// Caso 4: Deuda cero (no hay nada que pagar)
		{villa: 4, pago: 50_00, deuda: 0, expected: expected{destinado: 0, cuenta: 50_00}},
		// Caso 5: Pago cero (no se destina nada)
		{villa: 5, pago: 0, deuda: 100_00, expected: expected{destinado: 0, cuenta: 0}},
		// Caso 6: Pago mucho mayor que la deuda (500 USD vs 100 USD)
		{
			villa:    6,
			pago:     500_00,
			deuda:    100_00,
			expected: expected{destinado: 100_00, cuenta: 400_00},
		},
		// Caso 7: Pago muy pequeño frente a deuda grande (10 USD vs 1000 USD)
		{villa: 7, pago: 10_00, deuda: 1000_00, expected: expected{destinado: 10_00, cuenta: 0}},
		// Caso 8: Pago igual a deuda grande (1000 USD = 100000 centavos)
		{
			villa:    8,
			pago:     1000_00,
			deuda:    1000_00,
			expected: expected{destinado: 1000_00, cuenta: 0},
		},
	}

	for _, test := range tests {
		input := test
		t.Run(fmt.Sprintf("Villa_%d", input.villa), func(t *testing.T) {
			t.Parallel()

			db := testingutils.SetupInMemoryTestDB(t)

			// Insertando Cuota
			if err := db.Create(&administracion.Cuota{
				ID:    "c0",
				Monto: input.deuda,
			}).Error; err != nil {
				t.Fatalf("error insertando cuota: %v", err)
			}

			// Insertando Deuda
			if err := db.Create(&villas.IDeuda{
				ID:    "d0",
				Villa: input.villa,
				Cuota: "c0",
			}).Error; err != nil {
				t.Fatalf("error insertando deuda: %v", err)
			}

			// Insertando Pago
			if err := db.Create(&pagos.IPago{
				ID:     "p0",
				Villa:  input.villa,
				Metodo: metododepago.Efectivo,
				Monto:  input.pago / 100 * 325_42,
				Moneda: moneda.VED,
				Tasa:   325_42,
			}).Error; err != nil {
				t.Fatalf("error insertando pago: %v", err)
			}

			ctx := context.New(t.Context(), db)

			uc := command.NewRegistarPagoADeuda()

			dto := command.RegistarPagoADeudaDTO{PagoID: "p0"}

			_, cerr := uc.Exec(ctx, dto)

			require.NoError(t, cerr)

			var pago pagos.Pago
			err := db.First(&pago, "id = ?", dto.PagoID).Error

			if err != nil {
				t.Fatalf("error obteniendo pago: %v", err)
			}

			var destinado int
			if err := db.Model(&pagos.DestinoDePago{}).Where("pago = ?", "p0").Pluck("destinado", &destinado).Error; err != nil {
				t.Fatalf("error obteniendo  monto destinado: %v", err)
			}

			assert.Equal(
				t,
				input.expected.destinado,
				destinado,
				"El monto destinado para un pago de '%d USD' sobre una deuda de '%d USD' deberia ser de '%d USD'",
				input.pago/100,
				input.deuda/100,
				input.expected.destinado/100,
			)

			assert.Equal(
				t,
				input.expected.cuenta,
				pago.Cuenta,
				"El pago de '%d USD' debio resultar en una cuenta de '%d USD' tras pagar una deuda de '%d USD'",
				input.pago/100,
				input.expected.cuenta/100,
				input.deuda/100,
			)

		})
	}
}

func TestRegistarPagoADeuda_DeudaMasAntigua(t *testing.T) {
	db := testingutils.SetupInMemoryTestDB(t)

	if err := db.Create(&[]administracion.Cuota{
		{
			ID:    "c0",
			Monto: 1,
		},
		{
			ID:    "c1",
			Monto: 1,
		},
	}).Error; err != nil {
		t.Fatalf("error insertando cuota: %v", err)
	}

	if err := db.Create(&[]villas.IDeuda{
		{
			ID:       "d0",
			Villa:    1,
			Cuota:    "c0",
			Registro: time.Date(2025, time.January, 1, 13, 0, 0, 0, time.UTC),
		},
		{
			ID:       "d1",
			Villa:    1,
			Cuota:    "c1",
			Registro: time.Date(2025, time.February, 1, 13, 0, 0, 0, time.UTC),
		},
	}).Error; err != nil {
		t.Fatalf("error insertando pago: %v", err)
	}

	if err := db.Create(&pagos.IPago{
		ID:     "p0",
		Villa:  1,
		Monto:  100,
		Moneda: moneda.USD,
	}).Error; err != nil {
		t.Fatalf("error insertando pago: %v", err)
	}

	ctx := context.New(t.Context(), db)

	uc := command.NewRegistarPagoADeuda()

	dto := command.RegistarPagoADeudaDTO{"p0"}

	_, cerr := uc.Exec(ctx, dto)

	require.NoError(t, cerr)

	var deuda_destino string
	if err := db.Model(&pagos.DestinoDePago{}).Where("pago = ?", "p0").Pluck("deuda", &deuda_destino).Error; err != nil {
		t.Fatalf("error obteniendo deuda destino: %v", err)
	}

	assert.Equal(
		t,
		"d0",
		deuda_destino,
		"La deuda destino deberia ser 'd0' debido a que es la mas antigua",
	)
}
