package transaccion_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	currency "github.com/Sanaruca/condominio/internal/core/common/moneda"
	"github.com/Sanaruca/condominio/internal/core/common/quantity"
	"github.com/Sanaruca/condominio/internal/finanzas/models/operacion"
	"github.com/Sanaruca/condominio/internal/finanzas/models/transaccion"
	"github.com/Sanaruca/condominio/internal/finanzas/types/metodoperacion"
)

func nuevaOperacion(
	t *testing.T,
	factory *operacion.OperacionFactory,
	unidad string,
	moneda currency.Moneda,
	fecha time.Time,
) operacion.Operacion {
	t.Helper()
	o, err := factory.NuevoPago(
		unidad,
		"Pago de cuota",
		quantity.New(50000, 2),
		moneda,
		metodoperacion.Efectivo,
		quantity.New(10000, 2),
		fecha,
		"user-1",
		"",
		"",
	)
	require.NoError(t, err)
	return *o
}

func TestNuevo(t *testing.T) {
	factory := transaccion.NewTransaccionFactory()
	operacionFactory := operacion.NewOperacionFactory()
	fecha := time.Date(2026, 8, 15, 10, 0, 0, 0, time.UTC)

	t.Run("crea una transaccion con dos operaciones", func(t *testing.T) {
		ops := []operacion.Operacion{
			nuevaOperacion(t, operacionFactory, "U-001", currency.USD, fecha),
			nuevaOperacion(t, operacionFactory, "U-002", currency.USD, fecha),
		}
		tx, err := factory.Nuevo(fecha, "Compensacion con proveedor", "user-1", ops)
		require.NoError(t, err)
		require.NotNil(t, tx)

		assert.NotEmpty(t, tx.ID())
		assert.Equal(t, fecha, tx.Fecha())
		assert.Equal(t, "Compensacion con proveedor", tx.Concepto())
		assert.Equal(t, "user-1", tx.RegistradoPor())
		assert.Len(t, tx.Operaciones(), 2)
		assert.NotEmpty(t, tx.Registro())
	})

	t.Run("rechaza menos de dos operaciones", func(t *testing.T) {
		ops := []operacion.Operacion{
			nuevaOperacion(t, operacionFactory, "U-001", currency.USD, fecha),
		}
		_, err := factory.Nuevo(fecha, "Compensacion", "user-1", ops)
		assert.ErrorIs(t, err, transaccion.ErrOperacionesInsuficientes)
	})

	t.Run("rechaza operaciones con monedas distintas", func(t *testing.T) {
		ops := []operacion.Operacion{
			nuevaOperacion(t, operacionFactory, "U-001", currency.USD, fecha),
			nuevaOperacion(t, operacionFactory, "U-002", currency.VED, fecha),
		}
		_, err := factory.Nuevo(fecha, "Compensacion", "user-1", ops)
		assert.ErrorIs(t, err, transaccion.ErrMonedasDistintas)
	})
}

func TestAssemble(t *testing.T) {
	factory := transaccion.NewTransaccionFactory()
	operacionFactory := operacion.NewOperacionFactory()
	fecha := time.Date(2026, 8, 15, 10, 0, 0, 0, time.UTC)
	registro := time.Date(2026, 8, 15, 12, 30, 0, 0, time.UTC)

	ops := []operacion.Operacion{
		nuevaOperacion(t, operacionFactory, "U-001", currency.USD, fecha),
		nuevaOperacion(t, operacionFactory, "U-002", currency.USD, fecha),
	}

	tx := factory.Assemble("tx-1", fecha, "Compensacion", "user-1", registro, ops)
	require.NotNil(t, tx)

	assert.Equal(t, "tx-1", tx.ID())
	assert.Equal(t, fecha, tx.Fecha())
	assert.Equal(t, registro, tx.Registro())
	assert.Len(t, tx.Operaciones(), 2)
}
