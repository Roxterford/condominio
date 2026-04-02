package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/Sanaruca/condominio/graph"
	proveedorGorm "github.com/Sanaruca/condominio/internal/administracion/adapters/gorm"
	"github.com/Sanaruca/condominio/internal/administracion/models/gasto"
	"github.com/Sanaruca/condominio/internal/administracion/models/proveedor"
	administracionService "github.com/Sanaruca/condominio/internal/administracion/service"
	"github.com/Sanaruca/condominio/internal/core/common"
	coreContext "github.com/Sanaruca/condominio/internal/core/context"
	"github.com/Sanaruca/condominio/internal/core/envirotment"
	"github.com/Sanaruca/condominio/internal/core/session"
	pagosGorm "github.com/Sanaruca/condominio/internal/pagos/adapters/gorm"
	pagosRedis "github.com/Sanaruca/condominio/internal/pagos/adapters/redis"
	pagoService "github.com/Sanaruca/condominio/internal/pagos/service"
	"github.com/Sanaruca/condominio/internal/services/tasa"
	tasaCache "github.com/Sanaruca/condominio/internal/services/tasa/adapters/cache"
	tasaDolarAPI "github.com/Sanaruca/condominio/internal/services/tasa/adapters/dolarapi"
	tasaHybrid "github.com/Sanaruca/condominio/internal/services/tasa/adapters/hybrid"
	tasaLocal "github.com/Sanaruca/condominio/internal/services/tasa/adapters/local"
	sistemaService "github.com/Sanaruca/condominio/internal/sistema/service"
	"github.com/Sanaruca/condominio/internal/usuarios"
	usuariosGorm "github.com/Sanaruca/condominio/internal/usuarios/adapters/gorm"
	usuarioService "github.com/Sanaruca/condominio/internal/usuarios/service"
	villasGorm "github.com/Sanaruca/condominio/internal/villas/adapters/gorm"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"github.com/vektah/gqlparser/v2/ast"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const defaultPort = "8081"

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

	// Adapters / Dependencies
	eventBus := pagosRedis.NewRedisEventBus(redisClient, "pagos")

	// Repositories
	usuarioRepository := usuariosGorm.NewUsuarioGORMRepository(db, usuarios.NewFactory())
	villaRepository := villasGorm.NewVillaGORMRepository(db)
	pagoRepository := pagosGorm.NewPagoGORMRepository(db)
	proveedorRepository := proveedorGorm.NewGORMProveedorRepository(db, proveedorFactory)
	tasaLocalRepository := tasaLocal.NewGormLocalTasaRepository(db)
	tasaDolarAPIRepository := tasaDolarAPI.NewDolarAPITasaRepository()
	tasaRepository := tasaHybrid.NewHybridTasaRepository(
		tasaLocalRepository,
		[]tasa.TasaRepository{tasaDolarAPIRepository},
	)
	tasaCacheRepository := tasaCache.NewGormTasaCacheRepository(db)
	gastoRepository := proveedorGorm.NewGORMGastoRepository(db, gastoFactory)

	// Services
	tasaService := tasa.NewTasaService(tasaRepository, tasaCacheRepository)

	srv := handler.New(graph.NewExecutableSchema(graph.Config{Resolvers: graph.NewResolver(
		usuarioService.New(usuarioRepository),
		administracionService.New(
			proveedorRepository,
			gastoRepository,
			tasaService,
			proveedorFactory,
			emailFactory,
			phoneFactory,
		),
		pagoService.New(pagoRepository, villaRepository, eventBus, nil),
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

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}

func setupRedis() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
}

func setupDB() *gorm.DB {

	db, err := gorm.Open(sqlite.Open("dev.db"), &gorm.Config{
		Logger:         logger.Default.LogMode(logger.Info),
		TranslateError: true,
	})
	if err != nil {
		panic("failed to connect database")
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

}
