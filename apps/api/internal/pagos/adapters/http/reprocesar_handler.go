// Deprecated: Handler HTTP legacy de pagos. Usar transacciones en su lugar.
package http

import (
	"encoding/json"
	"net/http"

	"github.com/Sanaruca/condominio/internal/core/context"
	"github.com/Sanaruca/condominio/internal/pagos/app/command"
)

type ReprocesarHandler struct {
	reprocesarPagosHuerfanos command.ReprocesarPagosHuerfanos
}

func NewReprocesarHandler(
	reprocesarPagosHuerfanos command.ReprocesarPagosHuerfanos,
) *ReprocesarHandler {
	return &ReprocesarHandler{
		reprocesarPagosHuerfanos: reprocesarPagosHuerfanos,
	}
}

func (h *ReprocesarHandler) Reprocesar(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	baseCtx, err := context.Wrap(r.Context()).AsBase()
	if err != nil {
		http.Error(w, "Error al crear contexto", http.StatusInternalServerError)
		return
	}

	resultado, execErr := h.reprocesarPagosHuerfanos.Exec(baseCtx, nil)
	if execErr != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status": "error",
			"error":  execErr.Error(),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":    "success",
		"total":     resultado.Total,
		"procesados": resultado.Procesados,
		"errores":   resultado.Errores,
		"detalles":  resultado.Detalles,
	})
}
