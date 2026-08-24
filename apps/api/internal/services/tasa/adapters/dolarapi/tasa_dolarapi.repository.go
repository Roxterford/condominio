package dolarapi

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/Sanaruca/condominio/internal/core/common/quantity"
	"github.com/Sanaruca/condominio/internal/core/exception"
	"github.com/Sanaruca/condominio/internal/core/lib/logger"
	"github.com/Sanaruca/condominio/internal/services/tasa"
)

const baseURL = "https://ve.dolarapi.com/v1"

func convertirACentavos(valor float64) quantity.Quantity {
	return quantity.FromFloat(valor, quantity.DEFAULT_SCALE)
}

type DolarAPITasaRepository struct {
	client *http.Client
	cache  *tasaCache
	base   string
}

type tasaCache struct {
	ultimaConsulta time.Time
	tasa           tasa.Tasa
	ttl            time.Duration
	historial      map[string]tasa.Tasa // clave: fecha en formato YYYY-MM-DD
	historialTTL   time.Time
}

type dolarResponse struct {
	Fuente             string  `json:"fuente"`
	Nombre             string  `json:"nombre"`
	Compra             float64 `json:"compra"`
	Venta              float64 `json:"venta"`
	Promedio           float64 `json:"promedio"`
	FechaActualizacion string  `json:"fechaActualizacion"`
}

type historicoResponse []struct {
	Fuente   string  `json:"fuente"`
	Compra   float64 `json:"compra"`
	Venta    float64 `json:"venta"`
	Promedio float64 `json:"promedio"`
	Fecha    string  `json:"fecha"`
}

func NewDolarAPITasaRepository() *DolarAPITasaRepository {
	return newDolarAPITasaRepository(baseURL)
}

func newDolarAPITasaRepository(base string) *DolarAPITasaRepository {
	logger.Debug("🏭 Creando nueva instancia de DolarAPITasaRepository")
	return &DolarAPITasaRepository{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
		cache: &tasaCache{
			ttl:       5 * time.Minute,
			historial: make(map[string]tasa.Tasa),
		},
		base: base,
	}
}

func NewDolarAPITasaRepositoryWithURL(baseURL string) *DolarAPITasaRepository {
	return newDolarAPITasaRepository(baseURL)
}

func (a *DolarAPITasaRepository) Nombre() string {
	logger.Debug("📛 Obteniendo nombre del repositorio: dolarapi")
	return "dolarapi"
}

func (a *DolarAPITasaRepository) ObtenerTasaActual(tipo tasa.TipoDeCambio) (tasa.Tasa, error) {
	logger.Debug("💱 Iniciando obtención de tasa actual para tipo: %v", tipo)

	if a.cache.esValido() {
		logger.Debug("💾 Tasa encontrada en caché, valor: %d", a.cache.tasa.Valor.Value())
		return a.cache.tasa, nil
	}

	logger.Debug("🔄 Caché inválido, consultando API externa")
	_tasa, err := a.obtenerDesdeDolarAPI(tipo)
	if err != nil {
		logger.Debug("❌ Error al obtener tasa desde API: %v", err)
		return tasa.Tasa{}, err
	}

	a.cache.actualizar(_tasa)
	logger.Debug("✅ Tasa actualizada en caché, valor: %d", _tasa.Valor.Value())
	return _tasa, nil
}

func (a *DolarAPITasaRepository) ObtenerTasaPorFecha(
	tipo tasa.TipoDeCambio,
	fecha time.Time,
) (tasa.Tasa, error) {
	logger.Debug(
		"📅 Iniciando obtención de tasa por fecha - Tipo: %v, Fecha: %s",
		tipo,
		fecha.Format("2006-01-02"),
	)

	loc := time.Local
	hoy := time.Now().In(loc)
	hoy = time.Date(hoy.Year(), hoy.Month(), hoy.Day(), 0, 0, 0, 0, loc)
	fechaBusqueda := time.Date(fecha.Year(), fecha.Month(), fecha.Day(), 0, 0, 0, 0, loc)

	if fechaBusqueda.After(hoy) {
		logger.Debug("⏰ Fecha solicitada es futura, retornando error")
		return tasa.Tasa{}, tasa.ErrTasaNoDisponible
	}

	_tasa, err := a.obtenerHistorico(tipo, fechaBusqueda)
	if err != nil {
		logger.Debug("❌ Error al obtener tasa histórica: %v", err)
		return tasa.Tasa{}, err
	}

	logger.Debug("✅ Tasa histórica obtenida exitosamente, valor: %d", _tasa.Valor.Value())
	return _tasa, nil
}

func (a *DolarAPITasaRepository) ObtenerHistorico(
	tipo tasa.TipoDeCambio,
	desde, hasta time.Time,
) ([]tasa.Tasa, error) {
	logger.Debug(
		"📊 Iniciando obtención de tasas históricas - Tipo: %v, Desde: %s, Hasta: %s",
		tipo,
		desde.Format("2006-01-02"),
		hasta.Format("2006-01-02"),
	)

	endpoint := a.GetEndpoint(tipo)
	url := a.base + "/historicos" + endpoint
	logger.Debug("🌐 URL de consulta histórica: %s", url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logger.Debug("❌ Error al crear request HTTP: %v", err)
		return nil, exception.Wrap(err)
	}

	req.Header.Set("Accept", "application/json")

	resp, err := a.client.Do(req)
	if err != nil {
		logger.Debug("❌ Error al realizar request HTTP: %v", err)
		return nil, tasa.ErrSinConexion
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logger.Debug("⚠️ Respuesta HTTP con código: %d", resp.StatusCode)
		return nil, tasa.ErrTasaNoDisponible
	}

	var results historicoResponse
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		logger.Debug("❌ Error al decodificar respuesta JSON: %v", err)
		return nil, exception.Wrap(err)
	}

	logger.Debug("📋 Se obtuvieron %d registros históricos", len(results))
	tasas := make([]tasa.Tasa, 0, len(results))
	for _, r := range results {
		parsedFecha, _ := time.Parse("2006-01-02", r.Fecha)
		if parsedFecha.Before(desde) || parsedFecha.After(hasta) {
			continue
		}

		valor := convertirACentavos(r.Promedio)
		if valor.Value() == 0 {
			valor = convertirACentavos(r.Venta)
		}

		tasas = append(tasas, tasa.Tasa{
			Valor:  valor,
			Fuente: r.Fuente,
			Fecha:  parsedFecha,
			Tipo:   tipo,
			Moneda: "VES/USD",
		})
	}

	logger.Debug("✅ Se procesaron %d tasas en el rango solicitado", len(tasas))
	return tasas, nil
}

func (a *DolarAPITasaRepository) GetEndpoint(tipo tasa.TipoDeCambio) string {
	logger.Debug("🎯 Obteniendo endpoint para tipo de cambio: %v", tipo)

	switch tipo {
	case tasa.CambioParalelo:
		logger.Debug("🔗 Tipo paralelo, usando endpoint: /dolares/paralelo")
		return "/dolares/paralelo"
	case tasa.CambioPromedio:
		logger.Debug("🔗 Tipo promedio, usando endpoint: /dolares/oficial")
		return "/dolares/oficial"
	default:
		logger.Debug("🔗 Tipo no reconocido, usando endpoint por defecto: /dolares/oficial")
		return "/dolares/oficial"
	}
}

func (a *DolarAPITasaRepository) obtenerDesdeDolarAPI(tipo tasa.TipoDeCambio) (tasa.Tasa, error) {
	logger.Debug("🌐 Iniciando consulta de tasa actual a DolarAPI para tipo: %v", tipo)

	endpoint := a.GetEndpoint(tipo)
	url := a.base + endpoint

	logger.Debug("🔍 Consultando URL: %s", url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logger.Debug("❌ Error al crear request HTTP: %v", err)
		return tasa.Tasa{}, exception.Wrap(err)
	}

	req.Header.Set("Accept", "application/json")

	resp, err := a.client.Do(req)
	if err != nil {
		logger.Debug("❌ Error al realizar request HTTP: %v", err)
		return tasa.Tasa{}, tasa.ErrSinConexion
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logger.Debug("⚠️ Respuesta HTTP con código: %d", resp.StatusCode)
		return tasa.Tasa{}, tasa.ErrSinConexion
	}

	var result dolarResponse
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Debug("❌ Error al leer cuerpo de respuesta: %v", err)
		return tasa.Tasa{}, exception.Wrap(err)
	}

	logger.Debug("📨 Respuesta recibida: %s", string(body))

	if err := json.Unmarshal(body, &result); err != nil {
		logger.Debug("❌ Error al decodificar JSON: %v", err)
		return tasa.Tasa{}, exception.Wrap(err)
	}

	valor := convertirACentavos(result.Promedio)
	if valor.Value() == 0 {
		valor = convertirACentavos(result.Venta)
	}

	fecha, _ := time.Parse(time.RFC3339, result.FechaActualizacion)
	if fecha.IsZero() {
		fecha = time.Now()
	}

	logger.Debug(
		"💰 Tasa procesada - Valor: %d, Fuente: %s, Fecha: %s",
		valor.Value(),
		result.Fuente,
		fecha.Format("2006-01-02 15:04:05"),
	)

	return tasa.Tasa{
		Valor:  valor,
		Fuente: result.Fuente,
		Fecha:  fecha,
		Tipo:   tipo,
		Moneda: "VES/USD",
	}, nil
}

func (a *DolarAPITasaRepository) obtenerHistorico(
	tipo tasa.TipoDeCambio,
	fecha time.Time,
) (tasa.Tasa, error) {
	logger.Debug(
		"📅 Iniciando consulta de tasa histórica - Tipo: %v, Fecha: %s",
		tipo,
		fecha.Format("2006-01-02"),
	)

	// Primero intentar obtener desde el caché de historial
	if tasa, encontrada := a.obtenerTasaDesdeHistorial(fecha, tipo); encontrada {
		logger.Debug("✅ Tasa obtenida desde caché de historial")
		return tasa, nil
	}

	logger.Debug("🔄 Tasa no encontrada en caché, consultando API externa")

	endpoint := a.GetEndpoint(tipo)

	url := a.base + "/historicos" + endpoint
	logger.Debug("🔍 Consultando URL histórica: %s", url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logger.Debug("❌ Error al crear request HTTP: %v", err)
		return tasa.Tasa{}, exception.Wrap(err)
	}

	req.Header.Set("Accept", "application/json")

	resp, err := a.client.Do(req)
	if err != nil {
		logger.Debug("❌ Error al realizar request HTTP: %v", err)
		return tasa.Tasa{}, tasa.ErrSinConexion
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logger.Debug("⚠️ Respuesta HTTP con código: %d", resp.StatusCode)
		return tasa.Tasa{}, tasa.ErrTasaNoDisponible
	}

	var results historicoResponse
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		logger.Debug("❌ Error al decodificar respuesta JSON: %v", err)
		return tasa.Tasa{}, exception.Wrap(err)
	}

	logger.Debug("📋 Se obtuvieron %d registros históricos en total", len(results))

	if len(results) == 0 {
		logger.Debug("🔍 No se encontraron tasas históricas")
		return tasa.Tasa{}, tasa.ErrTasaNoEncontrada
	}

	// Buscar el registro que coincide con la fecha solicitada o la más cercana hacia atrás
	fechaBusqueda := fecha.Format("2006-01-02")
	var tasaEncontrada *struct {
		Fuente   string  `json:"fuente"`
		Compra   float64 `json:"compra"`
		Venta    float64 `json:"venta"`
		Promedio float64 `json:"promedio"`
		Fecha    string  `json:"fecha"`
	}
	var fechaEncontrada time.Time

	for i, r := range results {

		if r.Fecha == fechaBusqueda {
			logger.Debug("🎯 Tasa encontrada exacta en posición %d para fecha %s", i, fechaBusqueda)
			tasaEncontrada = &r
			fechaEncontrada, _ = time.Parse("2006-01-02", r.Fecha)
			break
		}
	}

	// Si no se encontró la fecha exacta, buscar la más cercana hacia atrás
	if tasaEncontrada == nil {
		logger.Debug(
			"🔍 No se encontró tasa exacta para %s, buscando fecha más cercana hacia atrás",
			fechaBusqueda,
		)

		fechaObjetivo, _ := time.Parse("2006-01-02", fechaBusqueda)
		var mejorDiferencia int = 999999 // Valor grande inicial

		for _, r := range results {
			fechaActual, err := time.Parse("2006-01-02", r.Fecha)
			if err != nil {
				continue
			}

			// Solo considerar fechas anteriores o iguales a la fecha buscada
			if fechaActual.After(fechaObjetivo) {
				continue
			}

			diferencia := int(fechaObjetivo.Sub(fechaActual).Hours() / 24)
			if diferencia < mejorDiferencia {
				mejorDiferencia = diferencia
				tasaEncontrada = &r
				fechaEncontrada = fechaActual
			}
		}

		if tasaEncontrada != nil {
			logger.Debug(
				"📅 Se encontró la tasa más cercana: %s (diferencia: %d días)",
				fechaEncontrada.Format("2006-01-02"),
				mejorDiferencia,
			)
		}
	}

	if tasaEncontrada == nil {
		logger.Debug("🔍 No se encontró ninguna tasa para la fecha solicitada: %s", fechaBusqueda)
		return tasa.Tasa{}, tasa.ErrTasaNoEncontrada
	}

	// Opcional: Guardar todo el historial en caché para futuras consultas
	logger.Debug("💾 Guardando historial completo en caché para optimización futura")
	a.guardarHistorialEnCache(results, tipo)

	r := *tasaEncontrada
	valor := convertirACentavos(r.Promedio)
	if valor.Value() == 0 {
		valor = convertirACentavos(r.Venta)
	}

	parsedFecha, _ := time.Parse("2006-01-02", r.Fecha)

	logger.Debug(
		"💰 Tasa histórica procesada - Valor: %d, Fuente: %s, Fecha: %s",
		valor.Value(),
		r.Fuente,
		parsedFecha.Format("2006-01-02"),
	)

	return tasa.Tasa{
		Valor:  valor,
		Fuente: r.Fuente,
		Fecha:  parsedFecha,
		Tipo:   tipo,
		Moneda: "VES/USD",
	}, nil
}

func (c *tasaCache) esValido() bool {
	if c.tasa.Valor.Value() == 0 {
		logger.Debug("🚫 Caché inválido: tasa con valor cero")
		return false
	}

	valido := time.Since(c.ultimaConsulta) < c.ttl
	logger.Debug(
		"🔍 Verificando validez de caché - Última consulta: %s, TTL: %v, Válido: %t",
		c.ultimaConsulta.Format("2006-01-02 15:04:05"),
		c.ttl,
		valido,
	)
	return valido
}

func (c *tasaCache) actualizar(tasa tasa.Tasa) {
	logger.Debug(
		"💾 Actualizando caché con nueva tasa - Valor: %d, Fuente: %s",
		tasa.Valor.Value(),
		tasa.Fuente,
	)
	c.tasa = tasa
	c.ultimaConsulta = time.Now()
}

func (c *tasaCache) historialValido() bool {
	// El historial es válido por 1 hora
	return time.Since(c.historialTTL) < time.Hour
}

func (a *DolarAPITasaRepository) guardarHistorialEnCache(
	results historicoResponse,
	tipo tasa.TipoDeCambio,
) {
	logger.Debug("💾 Guardando %d registros en caché de historial", len(results))

	// Limpiar historial anterior
	a.cache.historial = make(map[string]tasa.Tasa)

	for _, r := range results {
		parsedFecha, err := time.Parse("2006-01-02", r.Fecha)
		if err != nil {
			continue
		}

		valor := convertirACentavos(r.Promedio)
		if valor.Value() == 0 {
			valor = convertirACentavos(r.Venta)
		}

		a.cache.historial[r.Fecha] = tasa.Tasa{
			Valor:  valor,
			Fuente: r.Fuente,
			Fecha:  parsedFecha,
			Tipo:   tipo,
			Moneda: "VES/USD",
		}
	}

	a.cache.historialTTL = time.Now()
	logger.Debug("✅ Historial guardado en caché con %d registros", len(a.cache.historial))
}

func (a *DolarAPITasaRepository) obtenerTasaDesdeHistorial(
	fecha time.Time,
	tipo tasa.TipoDeCambio,
) (tasa.Tasa, bool) {
	fechaStr := fecha.Format("2006-01-02")

	if !a.cache.historialValido() {
		logger.Debug("🔍 Historial en caché inválido o expirado")
		return tasa.Tasa{}, false
	}

	if tasa, existe := a.cache.historial[fechaStr]; existe {
		logger.Debug("💾 Tasa encontrada en historial caché para fecha %s", fechaStr)
		return tasa, true
	}

	// Buscar fecha más cercana hacia atrás en el caché
	fechaObjetivo := fecha
	var mejorDiferencia int = 999999
	var tasaCercana tasa.Tasa
	var fechaCercana string

	for fechaStr, tasa := range a.cache.historial {
		if tasa.Tipo != tipo {
			continue
		}

		fechaCache, err := time.Parse("2006-01-02", fechaStr)
		if err != nil {
			continue
		}

		// Solo considerar fechas anteriores o iguales
		if fechaCache.After(fechaObjetivo) {
			continue
		}

		diferencia := int(fechaObjetivo.Sub(fechaCache).Hours() / 24)
		if diferencia < mejorDiferencia {
			mejorDiferencia = diferencia
			tasaCercana = tasa
			fechaCercana = fechaStr
		}
	}

	if mejorDiferencia < 999999 {
		logger.Debug(
			"📅 Tasa más cercana encontrada en caché: %s (diferencia: %d días)",
			fechaCercana,
			mejorDiferencia,
		)
		return tasaCercana, true
	}

	logger.Debug("🔍 No se encontró tasa en historial caché para fecha %s", fechaStr)
	return tasa.Tasa{}, false
}

func (a *DolarAPITasaRepository) GetCache() *tasaCache {
	return a.cache
}

func (c *tasaCache) EsValido() bool {
	return c.esValido()
}

func (c *tasaCache) Actualizar(tasa tasa.Tasa) {
	c.actualizar(tasa)
}

func (c *tasaCache) SetTTL(duration time.Duration) {
	c.ttl = duration
}

func (c *tasaCache) GuardarEnHistorial(fechaStr string, t tasa.Tasa) {
	if c.historial == nil {
		c.historial = make(map[string]tasa.Tasa)
	}
	c.historial[fechaStr] = t
	c.historialTTL = time.Now()
}

func (a *DolarAPITasaRepository) SetClient(client *http.Client) {
	a.client = client
}
