package dolarapi_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Sanaruca/condominio/internal/services/tasa"
	"github.com/Sanaruca/condominio/internal/services/tasa/adapters/dolarapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNombre(t *testing.T) {
	repo := dolarapi.NewDolarAPITasaRepository()
	assert.Equal(t, "dolarapi", repo.Nombre())
}

func TestGetEndpoint(t *testing.T) {
	tests := []struct {
		name     string
		tipo     tasa.TipoDeCambio
		expected string
	}{
		{
			name:     "tipo paralelo",
			tipo:     tasa.CambioParalelo,
			expected: "/dolares/paralelo",
		},
		{
			name:     "tipo promedio",
			tipo:     tasa.CambioPromedio,
			expected: "/dolares/oficial",
		},
		{
			name:     "tipo oficial",
			tipo:     tasa.CambioOficial,
			expected: "/dolares/oficial",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := dolarapi.NewDolarAPITasaRepository()
			result := repo.GetEndpoint(tt.tipo)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestObtenerTasaActual_CacheValido(t *testing.T) {
	repo := dolarapi.NewDolarAPITasaRepository()

	repo.GetCache().Actualizar(tasa.Tasa{
		Valor:  36500,
		Fuente: "Test",
		Fecha:  time.Now(),
		Tipo:   tasa.CambioParalelo,
		Moneda: "VES/USD",
	})

	result, err := repo.ObtenerTasaActual(tasa.CambioParalelo)

	require.NoError(t, err)
	assert.Equal(t, 36500, result.Valor)
	assert.Equal(t, "Test", result.Fuente)
}

func TestObtenerTasaActual_APIExitosa(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"fuente":             "BCV",
			"nombre":             "Dólar Paralelo",
			"compra":             35.50,
			"venta":              36.50,
			"promedio":           36.00,
			"fechaActualizacion": "2024-01-15T10:30:00-04:00",
		})
	}))
	defer server.Close()

	repo := dolarapi.NewDolarAPITasaRepositoryWithURL(server.URL)
	repo.SetClient(server.Client())

	result, err := repo.ObtenerTasaActual(tasa.CambioParalelo)

	require.NoError(t, err)
	assert.Equal(t, 36000, result.Valor)
	assert.Equal(t, "BCV", result.Fuente)
}

func TestObtenerTasaActual_APIDevuelveCeroPromedio(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"fuente":             "BCV",
			"nombre":             "Dólar Paralelo",
			"compra":             35.50,
			"venta":              36.50,
			"promedio":           0,
			"fechaActualizacion": "2024-01-15T10:30:00-04:00",
		})
	}))
	defer server.Close()

	repo := dolarapi.NewDolarAPITasaRepositoryWithURL(server.URL)
	repo.SetClient(server.Client())

	result, err := repo.ObtenerTasaActual(tasa.CambioParalelo)

	require.NoError(t, err)
	assert.Equal(t, 36500, result.Valor)
}

func TestObtenerTasaActual_ErrorHTTP(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	repo := dolarapi.NewDolarAPITasaRepositoryWithURL(server.URL)
	repo.SetClient(server.Client())

	_, err := repo.ObtenerTasaActual(tasa.CambioParalelo)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "Sin conexión")
}

func TestObtenerTasaActual_ErrorRed(t *testing.T) {
	repo := dolarapi.NewDolarAPITasaRepositoryWithURL("http://localhost:9999")
	repo.SetClient(&http.Client{
		Timeout: 1 * time.Millisecond,
	})

	_, err := repo.ObtenerTasaActual(tasa.CambioParalelo)

	require.Error(t, err)
}

func TestObtenerTasaActual_ErrorJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("invalid json"))
	}))
	defer server.Close()

	repo := dolarapi.NewDolarAPITasaRepositoryWithURL(server.URL)
	repo.SetClient(server.Client())

	_, err := repo.ObtenerTasaActual(tasa.CambioParalelo)

	require.Error(t, err)
}

func TestObtenerTasaPorFecha_FechaFutura(t *testing.T) {
	repo := dolarapi.NewDolarAPITasaRepository()

	fechaFutura := time.Now().Add(24 * time.Hour)
	fechaFutura = time.Date(fechaFutura.Year(), fechaFutura.Month(), fechaFutura.Day(), 0, 0, 0, 0, time.UTC)

	_, err := repo.ObtenerTasaPorFecha(tasa.CambioParalelo, fechaFutura)

	require.Error(t, err)
	assert.Equal(t, tasa.ErrTasaNoDisponible, err)
}

func TestObtenerTasaPorFecha_Exacta(t *testing.T) {
	fechaExacta := "2024-01-15"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]map[string]interface{}{
			{
				"fuente":   "BCV",
				"compra":   35.50,
				"venta":    36.50,
				"promedio": 36.00,
				"fecha":    fechaExacta,
			},
		})
	}))
	defer server.Close()

	repo := dolarapi.NewDolarAPITasaRepositoryWithURL(server.URL)
	repo.SetClient(server.Client())

	fecha, _ := time.Parse("2006-01-02", fechaExacta)
	result, err := repo.ObtenerTasaPorFecha(tasa.CambioParalelo, fecha)

	require.NoError(t, err)
	assert.Equal(t, 36000, result.Valor)
	assert.Equal(t, "BCV", result.Fuente)
}

func TestObtenerTasaPorFecha_BuscaFechaCercana(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]map[string]interface{}{
			{"fuente": "BCV", "compra": 35.50, "venta": 36.50, "promedio": 36.00, "fecha": "2024-01-10"},
			{"fuente": "BCV", "compra": 35.00, "venta": 36.00, "promedio": 35.50, "fecha": "2024-01-12"},
			{"fuente": "BCV", "compra": 34.50, "venta": 35.50, "promedio": 35.00, "fecha": "2024-01-14"},
		})
	}))
	defer server.Close()

	repo := dolarapi.NewDolarAPITasaRepositoryWithURL(server.URL)
	repo.SetClient(server.Client())

	fechaBusqueda, _ := time.Parse("2006-01-02", "2024-01-15")
	result, err := repo.ObtenerTasaPorFecha(tasa.CambioParalelo, fechaBusqueda)

	require.NoError(t, err)
	assert.Equal(t, 35000, result.Valor)
	assert.Equal(t, "2024-01-14", result.Fecha.Format("2006-01-02"))
}

func TestObtenerTasaPorFecha_ConHistorialCache(t *testing.T) {
	repo := dolarapi.NewDolarAPITasaRepository()

	repo.GetCache().GuardarEnHistorial("2024-01-14", tasa.Tasa{
		Valor:  35000,
		Fuente: "BCV",
		Fecha:  time.Date(2024, 1, 14, 0, 0, 0, 0, time.UTC),
		Tipo:   tasa.CambioParalelo,
		Moneda: "VES/USD",
	})

	fechaBusqueda, _ := time.Parse("2006-01-02", "2024-01-15")
	result, err := repo.ObtenerTasaPorFecha(tasa.CambioParalelo, fechaBusqueda)

	require.NoError(t, err)
	assert.Equal(t, 35000, result.Valor)
	assert.Equal(t, "2024-01-14", result.Fecha.Format("2006-01-02"))
}

func TestObtenerTasaPorFecha_NoEncontrada(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]map[string]interface{}{})
	}))
	defer server.Close()

	repo := dolarapi.NewDolarAPITasaRepositoryWithURL(server.URL)
	repo.SetClient(server.Client())

	fecha, _ := time.Parse("2006-01-02", "2024-01-15")
	_, err := repo.ObtenerTasaPorFecha(tasa.CambioParalelo, fecha)

	require.Error(t, err)
	assert.Equal(t, tasa.ErrTasaNoEncontrada, err)
}

func TestObtenerTasaPorFecha_ErrorHTTP(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	repo := dolarapi.NewDolarAPITasaRepositoryWithURL(server.URL)
	repo.SetClient(server.Client())

	fecha, _ := time.Parse("2006-01-02", "2024-01-15")
	_, err := repo.ObtenerTasaPorFecha(tasa.CambioParalelo, fecha)

	require.Error(t, err)
}

func TestObtenerHistorico_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]map[string]interface{}{
			{"fuente": "BCV", "compra": 35.50, "venta": 36.50, "promedio": 36.00, "fecha": "2024-01-10"},
			{"fuente": "BCV", "compra": 35.00, "venta": 36.00, "promedio": 35.50, "fecha": "2024-01-12"},
			{"fuente": "BCV", "compra": 34.50, "venta": 35.50, "promedio": 35.00, "fecha": "2024-01-14"},
		})
	}))
	defer server.Close()

	repo := dolarapi.NewDolarAPITasaRepositoryWithURL(server.URL)
	repo.SetClient(server.Client())

	desde, _ := time.Parse("2006-01-02", "2024-01-10")
	hasta, _ := time.Parse("2006-01-02", "2024-01-15")

	result, err := repo.ObtenerHistorico(tasa.CambioParalelo, desde, hasta)

	require.NoError(t, err)
	require.Len(t, result, 3)
	assert.Equal(t, 36000, result[0].Valor)
	assert.Equal(t, 35500, result[1].Valor)
	assert.Equal(t, 35000, result[2].Valor)
}

func TestObtenerHistorico_FiltroPorRango(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]map[string]interface{}{
			{"fuente": "BCV", "compra": 35.50, "venta": 36.50, "promedio": 36.00, "fecha": "2024-01-10"},
			{"fuente": "BCV", "compra": 35.00, "venta": 36.00, "promedio": 35.50, "fecha": "2024-01-12"},
			{"fuente": "BCV", "compra": 34.50, "venta": 35.50, "promedio": 35.00, "fecha": "2024-01-14"},
		})
	}))
	defer server.Close()

	repo := dolarapi.NewDolarAPITasaRepositoryWithURL(server.URL)
	repo.SetClient(server.Client())

	desde, _ := time.Parse("2006-01-02", "2024-01-12")
	hasta, _ := time.Parse("2006-01-02", "2024-01-14")

	result, err := repo.ObtenerHistorico(tasa.CambioParalelo, desde, hasta)

	require.NoError(t, err)
	require.Len(t, result, 2)
	assert.Equal(t, "2024-01-12", result[0].Fecha.Format("2006-01-02"))
	assert.Equal(t, "2024-01-14", result[1].Fecha.Format("2006-01-02"))
}

func TestObtenerHistorico_ErrorHTTP(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	repo := dolarapi.NewDolarAPITasaRepositoryWithURL(server.URL)
	repo.SetClient(server.Client())

	desde, _ := time.Parse("2006-01-02", "2024-01-10")
	hasta, _ := time.Parse("2006-01-02", "2024-01-15")

	_, err := repo.ObtenerHistorico(tasa.CambioParalelo, desde, hasta)

	require.Error(t, err)
	assert.Equal(t, tasa.ErrTasaNoDisponible, err)
}

func TestObtenerHistorico_ErrorJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("invalid json"))
	}))
	defer server.Close()

	repo := dolarapi.NewDolarAPITasaRepositoryWithURL(server.URL)
	repo.SetClient(server.Client())

	desde, _ := time.Parse("2006-01-02", "2024-01-10")
	hasta, _ := time.Parse("2006-01-02", "2024-01-15")

	_, err := repo.ObtenerHistorico(tasa.CambioParalelo, desde, hasta)

	require.Error(t, err)
}

func TestObtenerHistorico_PromedioCeroUsaVenta(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]map[string]interface{}{
			{"fuente": "BCV", "compra": 35.50, "venta": 36.50, "promedio": 0, "fecha": "2024-01-15"},
		})
	}))
	defer server.Close()

	repo := dolarapi.NewDolarAPITasaRepositoryWithURL(server.URL)
	repo.SetClient(server.Client())

	desde, _ := time.Parse("2006-01-02", "2024-01-10")
	hasta, _ := time.Parse("2006-01-02", "2024-01-15")

	result, err := repo.ObtenerHistorico(tasa.CambioParalelo, desde, hasta)

	require.NoError(t, err)
	require.Len(t, result, 1)
	assert.Equal(t, 36500, result[0].Valor)
}

func TestObtenerTasaActual_FechaActualizacionParsing(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"fuente":             "BCV",
			"nombre":             "Dólar Oficial",
			"compra":             35.50,
			"venta":              36.50,
			"promedio":           36.00,
			"fechaActualizacion": "2024-01-15T10:30:00-04:00",
		})
	}))
	defer server.Close()

	repo := dolarapi.NewDolarAPITasaRepositoryWithURL(server.URL)
	repo.SetClient(server.Client())

	result, err := repo.ObtenerTasaActual(tasa.CambioPromedio)

	require.NoError(t, err)
	assert.Equal(t, 36000, result.Valor)
	assert.Equal(t, "BCV", result.Fuente)
	assert.False(t, result.Fecha.IsZero())
}

func TestObtenerTasaActual_FechaActualizacionVacia(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"fuente":             "BCV",
			"nombre":             "Dólar Oficial",
			"compra":             35.50,
			"venta":              36.50,
			"promedio":           36.00,
			"fechaActualizacion": "",
		})
	}))
	defer server.Close()

	repo := dolarapi.NewDolarAPITasaRepositoryWithURL(server.URL)
	repo.SetClient(server.Client())

	result, err := repo.ObtenerTasaActual(tasa.CambioPromedio)

	require.NoError(t, err)
	assert.Equal(t, 36000, result.Valor)
	assert.False(t, result.Fecha.IsZero())
}

func TestCacheValido(t *testing.T) {
	repo := dolarapi.NewDolarAPITasaRepository()

	assert.False(t, repo.GetCache().EsValido())

	repo.GetCache().Actualizar(tasa.Tasa{
		Valor:  36500,
		Fuente: "BCV",
		Fecha:  time.Now(),
		Tipo:   tasa.CambioParalelo,
		Moneda: "VES/USD",
	})

	assert.True(t, repo.GetCache().EsValido())
}

func TestCacheInvalidoConValorCero(t *testing.T) {
	repo := dolarapi.NewDolarAPITasaRepository()

	repo.GetCache().Actualizar(tasa.Tasa{
		Valor:  0,
		Fuente: "BCV",
		Fecha:  time.Now(),
		Tipo:   tasa.CambioParalelo,
		Moneda: "VES/USD",
	})

	assert.False(t, repo.GetCache().EsValido())
}

func TestCacheTTLExpira(t *testing.T) {
	repo := dolarapi.NewDolarAPITasaRepository()

	repo.GetCache().Actualizar(tasa.Tasa{
		Valor:  36500,
		Fuente: "BCV",
		Fecha:  time.Now(),
		Tipo:   tasa.CambioParalelo,
		Moneda: "VES/USD",
	})

	repo.GetCache().SetTTL(-1 * time.Second)

	assert.False(t, repo.GetCache().EsValido())
}
