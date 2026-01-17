package main

import (
	"context"
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
	coreContext "github.com/Sanaruca/condominio/internal/core/context"
	"github.com/Sanaruca/condominio/internal/core/envirotment"
	"github.com/Sanaruca/condominio/internal/core/session"
	pagoService "github.com/Sanaruca/condominio/internal/pagos/service"
	"github.com/Sanaruca/condominio/internal/usuarios"
	usuarioService "github.com/Sanaruca/condominio/internal/usuarios/service"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"github.com/vektah/gqlparser/v2/ast"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const defaultPort = "8081"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	srv := handler.New(graph.NewExecutableSchema(graph.Config{Resolvers: graph.NewResolver(
		pagoService.New(),
		usuarioService.New(),
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
		(authMiddleware(coreContext.New(context.Background(), setupDB())))(srv),
	)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}

func setupDB() *gorm.DB {

	db, err := gorm.Open(sqlite.Open("dev.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	return db
}

func authMiddleware(ctx coreContext.BaseContext) func(http.Handler) http.Handler {

	return func(next http.Handler) http.Handler {
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
						usuario := &usuarios.UsuarioPayload{
							ID:    claims["ueid"].(string),
							Email: claims["email"].(string),
						}
						s := session.New(usuario)
						ctx = ctx.WithSession(s)
					}

				}
			}

			r = r.WithContext(context.WithValue(r.Context(), coreContext.BASE_CONTEXT_KEY, ctx))
			next.ServeHTTP(w, r)
		})
	}
}

func init() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

}
