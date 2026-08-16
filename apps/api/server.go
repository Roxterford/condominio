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
	"github.com/glebarez/sqlite"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"github.com/vektah/gqlparser/v2/ast"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/Sanaruca/condominio/graph"
	administracionGORM "github.com/Sanaruca/condominio/internal/administracion/adapters/gorm"
	"github.com/Sanaruca/condominio/internal/administracion/models/cuota"
	"github.com/Sanaruca/condominio/internal/administracion/models/deuda"
	"github.com/Sanaruca/condominio/internal/administracion/models/proveedor"
	administracionService "github.com/Sanaruca/condominio/internal/administracion/service"
	"github.com/Sanaruca/condominio/internal/core/common"
	"github.com/Sanaruca/condominio/internal/core/common/quantity"
	coreContext "github.com/Sanaruca/condominio/internal/core/context"
	"github.com/Sanaruca/condominio/internal/core/envirotment"
	"github.com/Sanaruca/condominio/internal/core/session"
	transaccionesGorm "github.com/Sanaruca/condominio/internal/finanzas/adapters/gorm"
	"github.com/Sanaruca/condominio/internal/finanzas/models/operacion"
	transaccionService "github.com/Sanaruca/condominio/internal/finanzas/service"
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

	// Factories
	emailFactory := common.NewEmailFactory([]string{})
	phoneFactory := common.NewPhoneFactory([]string{"58"}, []string{})
	quantityFactory := quantity.NewFactory(defaultDecimalPlaces)
	cuotaFactory := cuota.NewCuotaFactory(cuota.NewProyectoFactory(), quantityFactory)
	unidadFactory := unidad.NewUnidadFactory(quantityFactory)
	sujetoFactory := sujeto.NewSujetoFactory(emailFactory, phoneFactory)
	proveedorFactory := proveedor.NewProveedorFactory(
		common.NewEmailFactory([]string{}),
		common.NewPhoneFactory([]string{}, []string{}),
	)
	operacionFactory := operacion.NewOperacionFactory()

	// Repositories
	usuarioRepository := usuariosGorm.NewUsuarioGORMRepository(db, usuarios.NewFactory())
	unidadRepository := unidadesGorm.NewGORMUnidadRepository(db, unidadFactory, sujetoFactory)
	sujetoRepository := unidadesGorm.NewSujetoRepository(db, sujetoFactory)
	proveedorRepository := administracionGORM.NewGORMProveedorRepository(db, proveedorFactory)
	deudaFactory := deuda.NewDeudaFactory(quantityFactory)
	deudaRepository := administracionGORM.NewGORMDeudaRepository(db, deudaFactory, quantityFactory)
	recaudacionFinder := administracionGORM.NewGROMRecaudacionFinder(db, quantityFactory)
	gastosFinder := transaccionesGorm.NewGORMGastoFinder(db, quantityFactory, proveedorFactory)
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
	cuotaRepository := administracionGORM.NewGORMCuotaRepository(db, cuotaFactory)
	operacionRepository := transaccionesGorm.NewGORMOperacionRepository(
		db,
		quantityFactory,
		operacionFactory,
	)

	// Services
	tasaService := tasa.NewTasaService(tasaRepository, tasaCacheRepository)

	administracionService := administracionService.New(
		proveedorRepository,
		cuotaRepository,
		recaudacionFinder,
		proveedorFactory,
		emailFactory,
		phoneFactory,
		cuotaFactory,
	)

	transaccionServiceInstance := transaccionService.New(
		operacionRepository,
		gastosFinder,
		operacionFactory,
		unidadRepository,
		quantityFactory,
	)

	srv := handler.New(graph.NewExecutableSchema(graph.Config{Resolvers: graph.NewResolver(
		usuarioService.New(usuarioRepository),
		administracionService,
		unidadesService.NewUnidadesService(
			unidadRepository,
			sujetoRepository,
			deudaRepository,
			unidadEstadisticasFinder,
		),
		sistemaService.New(tasaService),
		transaccionServiceInstance,
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
