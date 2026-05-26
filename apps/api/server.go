package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/doganarif/govisual"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"github.com/vektah/gqlparser/v2/ast"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/Sanaruca/condominio/graph"
	administracionGORM "github.com/Sanaruca/condominio/internal/administracion/adapters/gorm"
	"github.com/Sanaruca/condominio/internal/administracion/models/cuota"
	"github.com/Sanaruca/condominio/internal/administracion/models/deuda"
	"github.com/Sanaruca/condominio/internal/administracion/models/gasto"
	"github.com/Sanaruca/condominio/internal/administracion/models/proveedor"
	administracionService "github.com/Sanaruca/condominio/internal/administracion/service"
	"github.com/Sanaruca/condominio/internal/core/common"
	"github.com/Sanaruca/condominio/internal/core/common/quantity"
	coreContext "github.com/Sanaruca/condominio/internal/core/context"
	"github.com/Sanaruca/condominio/internal/core/envirotment"
	"github.com/Sanaruca/condominio/internal/core/session"
	pagosGorm "github.com/Sanaruca/condominio/internal/pagos/adapters/gorm"
	pagosHTTP "github.com/Sanaruca/condominio/internal/pagos/adapters/http"
	pagosRedis "github.com/Sanaruca/condominio/internal/pagos/adapters/redis"
	"github.com/Sanaruca/condominio/internal/pagos/app/command"
	pagoConfig "github.com/Sanaruca/condominio/internal/pagos/config"
	"github.com/Sanaruca/condominio/internal/pagos/models/pago"
	pagoService "github.com/Sanaruca/condominio/internal/pagos/service"
	"github.com/Sanaruca/condominio/internal/services/tasa"
	tasaCache "github.com/Sanaruca/condominio/internal/services/tasa/adapters/cache"
	tasaDolarAPI "github.com/Sanaruca/condominio/internal/services/tasa/adapters/dolarapi"
	tasaHybrid "github.com/Sanaruca/condominio/internal/services/tasa/adapters/hybrid"
	tasaLocal "github.com/Sanaruca/condominio/internal/services/tasa/adapters/local"
	sistemaService "github.com/Sanaruca/condominio/internal/sistema/service"
	unidadesGorm "github.com/Sanaruca/condominio/internal/unidades/adapters/gorm"
	"github.com/Sanaruca/condominio/internal/unidades/models/sujeto"
	"github.com/Sanaruca/condominio/internal/unidades/models/unidad"
	unidadesService "github.com/Sanaruca/condominio/internal/unidades/service"
	"github.com/Sanaruca/condominio/internal/usuarios"
	usuariosGorm "github.com/Sanaruca/condominio/internal/usuarios/adapters/gorm"
	usuarioService "github.com/Sanaruca/condominio/internal/usuarios/service"
)

const defaultPort = "8081"
const defaultDecimalPlaces = quantity.DEFAULT_SCALE

func main() {

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	db := setupDB()
	redisClient := setupRedis()

	// Factories
	proveedorFactory := proveedor.NewProveedorFactory(
		common.NewEmailFactory([]string{}),
		common.NewPhoneFactory([]string{}, []string{}),
	)
	gastoFactory := gasto.NewGastoFactory()
	emailFactory := common.NewEmailFactory([]string{})
	phoneFactory := common.NewPhoneFactory([]string{"58"}, []string{})
	quantityFactory := quantity.NewFactory(defaultDecimalPlaces)
	cuotaFactory := cuota.NewCuotaFactory(cuota.NewProyectoFactory(), quantityFactory)
	deudaFactory := deuda.NewDeudaFactory(quantityFactory)
	unidadFactory := unidad.NewUnidadFactory(quantityFactory)
	sujetoFactory := sujeto.NewSujetoFactory(emailFactory, phoneFactory)
	pagoFactory := pago.NewPagoFactory()

	// Adapters / Dependencies
	eventBus := pagosRedis.NewRedisEventBus(redisClient, "pagos")

	// Repositories
	usuarioRepository := usuariosGorm.NewUsuarioGORMRepository(db, usuarios.NewFactory())
	unidadRepository := unidadesGorm.NewGORMUnidadRepository(db, unidadFactory, sujetoFactory)
	sujetoRepository := unidadesGorm.NewSujetoRepository(db, sujetoFactory)
	pagoRepository := pagosGorm.NewGORMPagoRepository(db, quantityFactory, pagoFactory)
	proveedorRepository := administracionGORM.NewGORMProveedorRepository(db, proveedorFactory)
	deudaRepository := administracionGORM.NewGORMDeudaRepository(db, deudaFactory, quantityFactory)
	recaudacionFinder := administracionGORM.NewGROMRecaudacionFinder(db, quantityFactory)
	unidadEstadisticasFinder := unidadesGorm.NewGORMUnidadEstadisticasFinder(
		db,
		quantityFactory,
	)
	tasaLocalRepository := tasaLocal.NewGormLocalTasaRepository(db)
	tasaDolarAPIRepository := tasaDolarAPI.NewDolarAPITasaRepository()
	tasaRepository := tasaHybrid.NewHybridTasaRepository(
		tasaLocalRepository,
		[]tasa.TasaRepository{tasaDolarAPIRepository},
	)
	tasaCacheRepository := tasaCache.NewGormTasaCacheRepository(db)
	gastoRepository := administracionGORM.NewGORMGastoRepository(db, gastoFactory, quantityFactory)
	cuotaRepository := administracionGORM.NewGORMCuotaRepository(db, cuotaFactory)

	// Services
	tasaService := tasa.NewTasaService(tasaRepository, tasaCacheRepository)

	// Configurar handlers de eventos para pagos
	aplicarPago := command.NewAplicarPago(
		pagoRepository,
		unidadRepository,
		deudaRepository,
	)

	eventHandlersConfig := pagoConfig.NewEventHandlersConfig(aplicarPago)

	// Handler HTTP para el worker
	workerHandler := pagosRedis.NewWorkerHTTPHandler(
		eventHandlersConfig,
		redisClient,
	)

	// Handler HTTP para reprocesar pagos huérfanos
	pagoServiceInstance := pagoService.New(
		pagoRepository,
		pagoFactory,
		unidadRepository,
		eventBus,
		tasaService,
		quantityFactory,
		aplicarPago,
	)
	reprocesarHandler := pagosHTTP.NewReprocesarHandler(
		pagoServiceInstance.Commands.ReprocesarPagosHuerfanos,
	)

	srv := handler.New(graph.NewExecutableSchema(graph.Config{Resolvers: graph.NewResolver(
		usuarioService.New(usuarioRepository),
		administracionService.New(
			proveedorRepository,
			gastoRepository,
			cuotaRepository,
			deudaRepository,
			recaudacionFinder,
			tasaService,
			proveedorFactory,
			emailFactory,
			phoneFactory,
			gastoFactory,
		),
		pagoServiceInstance,
		unidadesService.NewUnidadesService(
			unidadRepository,
			sujetoRepository,
			deudaRepository,
			unidadEstadisticasFinder,
		),
		sistemaService.New(tasaService),
	)}))

	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})

	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	mux := http.NewServeMux()

	mux.Handle("/", playground.Handler("GraphQL playground", "/query"))
	mux.Handle(
		"/query",
		(corsMiddleware)(authMiddleware(srv)),
	)

	// Endpoints para el worker de eventos
	mux.HandleFunc("/api/worker/process", workerHandler.ProcessEvents)
	mux.HandleFunc("/api/worker/health", workerHandler.HealthCheck)

	// Endpoint para reprocesar pagos huérfanos
	mux.HandleFunc("/api/pagos/reprocesar", reprocesarHandler.Reprocesar)

	var handler http.Handler = mux

	if envirotment.GetAppEnv() == envirotment.Dev {
		handler = govisual.Wrap(
			mux,
			govisual.WithRequestBodyLogging(true),
			govisual.WithResponseBodyLogging(true),
			govisual.WithIgnorePaths("/api/worker/health"),
		)
	}

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}

func setupRedis() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
}

func findProjectRoot() string {
	cwd, _ := os.Getwd()
	dir := cwd
	for i := 0; i < 10; i++ {
		if _, err := os.Stat(filepath.Join(dir, ".env")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return cwd
}

func setupDB() *gorm.DB {
	dsn := envirotment.Get(envirotment.DATABASE_URL)
	if dsn == "" {
		panic("DATABASE_URL is not set")
	}

	// Resolve relative SQLite paths against project root (where .env lives)
	if strings.HasPrefix(dsn, "file:./") {
		dsn = "file:" + filepath.Join(findProjectRoot(), dsn[6:])
	}

	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
		Logger:         logger.Default.LogMode(logger.Info),
		TranslateError: true,
	})
	if err != nil {
		panic("failed to connect database: " + err.Error())
	}

	return db
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Obtener hosts permitidos desde variables de entorno o usar defaults
		allowedHosts := strings.Split(os.Getenv("CORS_ALLOWED_HOSTS"), ",")
		if len(allowedHosts) == 1 && allowedHosts[0] == "" {
			// Si no hay configuración, usar hosts de desarrollo por defecto
			allowedHosts = []string{
				"http://localhost:3000",
				"http://localhost:3001",
				"http://127.0.0.1:3000",
				"http://127.0.0.1:3001",
				"http://localhost:4000",
				"http://localhost:4001",
				"http://127.0.0.1:4000",
				"http://127.0.0.1:4001",
			}
		}

		origin := r.Header.Get("Origin")

		// Verificar si el origin está en la lista de permitidos
		allowed := false
		for _, host := range allowedHosts {
			if strings.TrimSpace(host) == origin {
				allowed = true
				break
			}
		}

		if allowed {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		}

		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		// Manejar preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		authHeader := r.Header.Get("Authorization")
		if authHeader != "" {
			bearerToken := strings.Split(authHeader, " ")
			if len(bearerToken) == 2 {
				jwt_token := bearerToken[1]

				var claims jwt.MapClaims
				_, err := jwt.ParseWithClaims(
					jwt_token,
					&claims,
					func(t *jwt.Token) (any, error) {
						return []byte(envirotment.GetSecretKey()), nil
					},
				)

				if err == nil {
					usuario := &session.CredencialDeUsuario{
						ID:    claims["ueid"].(string),
						Email: claims["email"].(string),
					}
					ctx := coreContext.InjectUser(r.Context(), usuario)
					r = r.WithContext(ctx)
				}

			}
		}

		next.ServeHTTP(w, r)
	})
}

func init() {
	godotenv.Load()
	godotenv.Load("../../.env")
	godotenv.Load("../.env")
}
