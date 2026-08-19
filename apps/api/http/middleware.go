package httprouter

import (
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"

	coreContext "github.com/Sanaruca/condominio/internal/core/context"
	"github.com/Sanaruca/condominio/internal/core/envirotment"
	"github.com/Sanaruca/condominio/internal/core/session"
)

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

// correlationMiddleware extrae o genera un correlation ID y lo inyecta en el context.
func correlationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		correlationID := r.Header.Get("X-Correlation-ID")
		if correlationID == "" {
			correlationID = r.Header.Get("X-Request-ID")
		}
		ctx := coreContext.InjectCorrelationID(r.Context(), correlationID)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
