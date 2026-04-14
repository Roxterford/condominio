package redis

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Sanaruca/condominio/internal/pagos/config"
	"github.com/redis/go-redis/v9"
)

// WorkerHTTPHandler maneja las peticiones HTTP para procesar eventos
type WorkerHTTPHandler struct {
	eventConfig *config.EventHandlersConfig
	redisClient *redis.Client
}

// NewWorkerHTTPHandler crea una nueva instancia del handler HTTP
func NewWorkerHTTPHandler(
	eventConfig *config.EventHandlersConfig,
	redisClient *redis.Client,
) *WorkerHTTPHandler {
	return &WorkerHTTPHandler{
		eventConfig: eventConfig,
		redisClient: redisClient,
	}
}

// ProcessEvents procesa todos los eventos pendientes en Redis via HTTP
func (h *WorkerHTTPHandler) ProcessEvents(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Crear worker con el dispatcher configurado
	worker := NewRedisEventWorker(
		h.redisClient,
		"pagos",           // stream name
		"pagos-group",     // consumer group
		"cron-worker",     // consumer name
		h.eventConfig.Dispatcher,
	)

	// Ejecutar worker en el contexto de la petición
	err := worker.Start(r.Context())
	if err != nil {
		log.Printf("Error procesando eventos: %v", err)
		http.Error(w, "Error procesando eventos", http.StatusInternalServerError)
		return
	}

	// Responder con éxito
	w.Header().Set("Content-Type", "application/json")
	response := map[string]string{
		"status":  "success",
		"message": "Eventos procesados correctamente",
	}
	json.NewEncoder(w).Encode(response)
}

// HealthCheck verifica el estado del worker via HTTP
func (h *WorkerHTTPHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Verificar conexión a Redis
	_, err := h.redisClient.Ping(r.Context()).Result()
	if err != nil {
		http.Error(w, "Redis no disponible", http.StatusServiceUnavailable)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"status": "healthy",
		"redis":  "connected",
		"stream": "pagos",
	}
	json.NewEncoder(w).Encode(response)
}
