package operacion_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	currency "github.com/Sanaruca/condominio/internal/core/common/moneda"
	"github.com/Sanaruca/condominio/internal/core/common/quantity"
	"github.com/Sanaruca/condominio/internal/finanzas/models/operacion"
	"github.com/Sanaruca/condominio/internal/finanzas/types/metodoperacion"
	"github.com/Sanaruca/condominio/internal/finanzas/types/roldestionoperacion"
	"github.com/Sanaruca/condominio/internal/finanzas/types/tipoperacion"
)

func TestNuevo(t *testing.T) {
	factory := operacion.NewOperacionFactory()
	fecha := time.Date(2026, 8, 15, 10, 0, 0, 0, time.UTC)
	unidad := "U-001"

	t.Run("crea una operacion valida", func(t *testing.T) {
		o, err := factory.Nuevo(
			"Pago de cuota",
			quantity.New(50000, 2),
			currency.USD,
			metodoperacion.Efectivo,
			quantity.New(10000, 2),
			tipoperacion.Credito,
			roldestionoperacion.Unidad,
			nil,
			&unidad,
			nil,
			fecha,
			"user-1",
		)
		require.NoError(t, err)
		require.NotNil(t, o)

		assert.NotEmpty(t, o.ID())
		assert.Equal(t, fecha, o.Fecha())
		assert.Equal(t, "Pago de cuota", o.Concepto())
		assert.Equal(t, int64(50000), o.Monto().Value())
		assert.Equal(t, currency.USD, o.Moneda())
		assert.Equal(t, metodoperacion.Efectivo, o.Metodo())
		assert.Equal(t, tipoperacion.Credito, o.Tipo())
		assert.Equal(t, roldestionoperacion.Unidad, o.Rol())
		assert.Equal(t, "U-001", *o.UnidadCodigo())
		assert.Nil(t, o.ProveedorID())
		assert.Nil(t, o.CuotaID())
		assert.Equal(t, "user-1", o.RegistradoPor())
		assert.True(t, o.EsCredito())
		assert.False(t, o.EsDebito())
		assert.NotEmpty(t, o.PullEvents())
	})

	t.Run("rechaza monto menor o igual a cero", func(t *testing.T) {
		_, err := factory.Nuevo(
			"Pago de cuota",
			quantity.New(0, 2),
			currency.USD,
			metodoperacion.Efectivo,
			quantity.New(10000, 2),
			tipoperacion.Credito,
			roldestionoperacion.Unidad,
			nil,
			&unidad,
			nil,
			fecha,
			"user-1",
		)
		assert.ErrorIs(t, err, operacion.ErrMontoInvalido)
	})

	t.Run("rechaza moneda invalida", func(t *testing.T) {
		_, err := factory.Nuevo(
			"Pago de cuota",
			quantity.New(50000, 2),
			currency.Moneda("XXX"),
			metodoperacion.Efectivo,
			quantity.New(10000, 2),
			tipoperacion.Credito,
			roldestionoperacion.Unidad,
			nil,
			&unidad,
			nil,
			fecha,
			"user-1",
		)
		assert.Error(t, err)
	})

	t.Run("rechaza metodo invalido", func(t *testing.T) {
		_, err := factory.Nuevo(
			"Pago de cuota",
			quantity.New(50000, 2),
			currency.USD,
			metodoperacion.MetodoDeOperacion("X"),
			quantity.New(10000, 2),
			tipoperacion.Credito,
			roldestionoperacion.Unidad,
			nil,
			&unidad,
			nil,
			fecha,
			"user-1",
		)
		assert.Error(t, err)
	})

	t.Run("rechaza tipo invalido", func(t *testing.T) {
		_, err := factory.Nuevo(
			"Pago de cuota",
			quantity.New(50000, 2),
			currency.USD,
			metodoperacion.Efectivo,
			quantity.New(10000, 2),
			tipoperacion.TipoDeOperacion("X"),
			roldestionoperacion.Unidad,
			nil,
			&unidad,
			nil,
			fecha,
			"user-1",
		)
		assert.Error(t, err)
	})

	t.Run("rechaza rol invalido", func(t *testing.T) {
		_, err := factory.Nuevo(
			"Pago de cuota",
			quantity.New(50000, 2),
			currency.USD,
			metodoperacion.Efectivo,
			quantity.New(10000, 2),
			tipoperacion.Credito,
			roldestionoperacion.RolDestinoDeOperacion("X"),
			nil,
			&unidad,
			nil,
			fecha,
			"user-1",
		)
		assert.ErrorIs(t, err, operacion.ErrRolInvalido)
	})

	t.Run("exige unidad cuando el rol es UNIDAD", func(t *testing.T) {
		_, err := factory.Nuevo(
			"Pago de cuota",
			quantity.New(50000, 2),
			currency.USD,
			metodoperacion.Efectivo,
			quantity.New(10000, 2),
			tipoperacion.Credito,
			roldestionoperacion.Unidad,
			nil,
			nil,
			nil,
			fecha,
			"user-1",
		)
		assert.ErrorIs(t, err, operacion.ErrUnidadRequerida)
	})

	t.Run("exige proveedor cuando el rol es PROVEEDOR", func(t *testing.T) {
		_, err := factory.Nuevo(
			"Gasto de mantenimiento",
			quantity.New(50000, 2),
			currency.USD,
			metodoperacion.TransferenciaNacional,
			quantity.New(10000, 2),
			tipoperacion.Debito,
			roldestionoperacion.Proveedor,
			nil,
			nil,
			nil,
			fecha,
			"user-1",
		)
		assert.ErrorIs(t, err, operacion.ErrProveedorRequerido)
	})

	t.Run("permite rol CONDOMINIO sin referencias", func(t *testing.T) {
		o, err := factory.Nuevo(
			"Gasto del condominio",
			quantity.New(50000, 2),
			currency.USD,
			metodoperacion.Compensacion,
			quantity.New(10000, 2),
			tipoperacion.Debito,
			roldestionoperacion.Condominio,
			nil,
			nil,
			nil,
			fecha,
			"user-1",
		)
		require.NoError(t, err)
		assert.Equal(t, roldestionoperacion.Condominio, o.Rol())
		assert.Nil(t, o.UnidadCodigo())
		assert.Nil(t, o.ProveedorID())
	})
}

func TestNuevoPago(t *testing.T) {
	factory := operacion.NewOperacionFactory()
	fecha := time.Date(2026, 8, 15, 10, 0, 0, 0, time.UTC)

	o, err := factory.NuevoPago(
		"U-001",
		"Pago de cuota",
		quantity.New(50000, 2),
		currency.VED,
		metodoperacion.PagoMovil,
		quantity.New(360000, 2),
		fecha,
		"user-1",
	)
	require.NoError(t, err)

	assert.Equal(t, tipoperacion.Credito, o.Tipo())
	assert.Equal(t, roldestionoperacion.Unidad, o.Rol())
	assert.Equal(t, "U-001", *o.UnidadCodigo())
	assert.Equal(t, currency.VED, o.Moneda())
	assert.Equal(t, int64(50000), o.Monto().Value())
}

func TestNuevoGasto(t *testing.T) {
	factory := operacion.NewOperacionFactory()
	fecha := time.Date(2026, 8, 15, 10, 0, 0, 0, time.UTC)

	t.Run("gasto a proveedor con cuota", func(t *testing.T) {
		proveedor := "prov-1"
		cuota := "cuota-1"
		o, err := factory.NuevoGasto(
			"Mantenimiento de ascensor",
			&proveedor,
			false,
			quantity.New(50000, 2),
			currency.USD,
			metodoperacion.TransferenciaNacional,
			quantity.New(10000, 2),
			fecha,
			&cuota,
			"user-1",
		)
		require.NoError(t, err)

		assert.Equal(t, tipoperacion.Debito, o.Tipo())
		assert.Equal(t, roldestionoperacion.Proveedor, o.Rol())
		assert.Equal(t, "prov-1", *o.ProveedorID())
		assert.Equal(t, "cuota-1", *o.CuotaID())
	})

	t.Run("gasto al condominio", func(t *testing.T) {
		o, err := factory.NuevoGasto(
			"Pago de servicio comun",
			nil,
			true,
			quantity.New(50000, 2),
			currency.USD,
			metodoperacion.Cheque,
			quantity.New(10000, 2),
			fecha,
			nil,
			"user-1",
		)
		require.NoError(t, err)

		assert.Equal(t, roldestionoperacion.Condominio, o.Rol())
		assert.Nil(t, o.ProveedorID())
	})

	t.Run("rechaza gasto a proveedor sin proveedor", func(t *testing.T) {
		_, err := factory.NuevoGasto(
			"Mantenimiento",
			nil,
			false,
			quantity.New(50000, 2),
			currency.USD,
			metodoperacion.Efectivo,
			quantity.New(10000, 2),
			fecha,
			nil,
			"user-1",
		)
		assert.ErrorIs(t, err, operacion.ErrProveedorRequerido)
	})
}

func TestNuevoReembolso(t *testing.T) {
	factory := operacion.NewOperacionFactory()
	fecha := time.Date(2026, 8, 15, 10, 0, 0, 0, time.UTC)

	o, err := factory.NuevoReembolso(
		"U-001",
		"Reembolso por pago duplicado",
		quantity.New(20000, 2),
		currency.USD,
		metodoperacion.Efectivo,
		quantity.New(10000, 2),
		fecha,
		"user-1",
	)
	require.NoError(t, err)

	assert.Equal(t, tipoperacion.Debito, o.Tipo())
	assert.Equal(t, roldestionoperacion.Unidad, o.Rol())
	assert.Equal(t, "U-001", *o.UnidadCodigo())
}

func TestAssemble(t *testing.T) {
	factory := operacion.NewOperacionFactory()
	fecha := time.Date(2026, 8, 15, 10, 0, 0, 0, time.UTC)
	registro := time.Date(2026, 8, 15, 12, 30, 0, 0, time.UTC)
	unidad := "U-001"

	o := factory.Assemble(
		"op-1",
		fecha,
		"Pago de cuota",
		quantity.New(50000, 2),
		currency.USD,
		metodoperacion.Efectivo,
		quantity.New(10000, 2),
		tipoperacion.Credito,
		roldestionoperacion.Unidad,
		nil,
		&unidad,
		nil,
		"user-1",
		registro,
	)
	require.NotNil(t, o)

	assert.Equal(t, "op-1", o.ID())
	assert.Equal(t, fecha, o.Fecha())
	assert.Equal(t, registro, o.Registro())
	assert.Equal(t, int64(50000), o.Monto().Value())
}

func TestMontoEnUSD(t *testing.T) {
	factory := operacion.NewOperacionFactory()
	fecha := time.Date(2026, 8, 15, 10, 0, 0, 0, time.UTC)

	t.Run("moneda USD sin conversion", func(t *testing.T) {
		o, err := factory.NuevoPago(
			"U-001",
			"Pago",
			quantity.New(50000, 2),
			currency.USD,
			metodoperacion.Efectivo,
			quantity.New(10000, 2),
			fecha,
			"user-1",
		)
		require.NoError(t, err)
		assert.Equal(t, int64(50000), o.MontoEnUSD().Value())
	})

	t.Run("moneda VED divide por la tasa", func(t *testing.T) {
		o, err := factory.NuevoPago(
			"U-001",
			"Pago",
			quantity.New(50000, 2),
			currency.VED,
			metodoperacion.PagoMovil,
			quantity.New(500, 2),
			fecha,
			"user-1",
		)
		require.NoError(t, err)
		assert.Equal(t, int64(10000), o.MontoEnUSD().Value())
	})
}
