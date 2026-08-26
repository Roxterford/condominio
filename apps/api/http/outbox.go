package httprouter

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	gormAdapter "github.com/Sanaruca/condominio/internal/core/adapters/gorm"
	appLogger "github.com/Sanaruca/condominio/internal/core/lib/logger"
)

// handleOutboxStats returns summary statistics for outbox events
func (s *Server) handleOutboxStats(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pending, err := s.outboxStore.CountPending(ctx)
	if err != nil {
		appLogger.ErrorCtx(ctx, err, "failed to count pending outbox events")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	// Also count DLQ entries
	var dlqCount int64
	if err := s.outboxStore.DB().Model(&gormAdapter.OutboxDLQ{}).Count(&dlqCount).Error; err != nil {
		appLogger.ErrorCtx(ctx, err, "failed to count DLQ entries")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"pending":   pending,
		"dlq_count": dlqCount,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// handleOutboxPending returns paginated list of pending outbox events
func (s *Server) handleOutboxPending(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	limit := 100
	if l := r.URL.Query().Get("limit"); l != "" {
		if v, err := strconv.Atoi(l); err == nil && v > 0 && v <= 500 {
			limit = v
		}
	}

	events, err := s.outboxStore.GetPending(ctx, limit)
	if err != nil {
		appLogger.ErrorCtx(ctx, err, "failed to get pending outbox events")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"events":    events,
		"count":     len(events),
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// handleOutboxDLQ returns paginated list of DLQ events
func (s *Server) handleOutboxDLQ(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	limit := 100
	if l := r.URL.Query().Get("limit"); l != "" {
		if v, err := strconv.Atoi(l); err == nil && v > 0 && v <= 500 {
			limit = v
		}
	}

	var dlqEvents []*gormAdapter.OutboxDLQ
	err := s.outboxStore.DB().
		Order("failed_at DESC").
		Limit(limit).
		Find(&dlqEvents).Error
	if err != nil {
		appLogger.ErrorCtx(ctx, err, "failed to get DLQ events")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"events":    dlqEvents,
		"count":     len(dlqEvents),
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// handleOutboxDlqRetry retries a specific DLQ event
func (s *Server) handleOutboxDlqRetry(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "missing id parameter", http.StatusBadRequest)
		return
	}

	// Move from DLQ back to outbox_events with retries reset
	tx := s.outboxStore.DB().Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var dlqEvent gormAdapter.OutboxDLQ
	if err := tx.Where("id = ?", id).First(&dlqEvent).Error; err != nil {
		appLogger.ErrorCtx(ctx, err, "DLQ event not found", "id", id)
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	// Create new outbox event from DLQ entry
	outboxEvent := &gormAdapter.OutboxEvent{
		ID:            dlqEvent.ID,
		EventName:     dlqEvent.EventName,
		Payload:       dlqEvent.Payload,
		CorrelationID: dlqEvent.CorrelationID,
		CausationID:   dlqEvent.CausationID,
		OccurredAt:    dlqEvent.OccurredAt,
		CreatedAt:     time.Now().UTC(),
		Retries:       0,
	}

	if err := tx.Create(outboxEvent).Error; err != nil {
		tx.Rollback()
		appLogger.ErrorCtx(ctx, err, "failed to re-create outbox event from DLQ", "id", id)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	if err := tx.Where("id = ?", id).Delete(&gormAdapter.OutboxDLQ{}).Error; err != nil {
		tx.Rollback()
		appLogger.ErrorCtx(ctx, err, "failed to delete DLQ entry", "id", id)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	if err := tx.Commit().Error; err != nil {
		appLogger.ErrorCtx(ctx, err, "failed to commit DLQ retry", "id", id)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	appLogger.InfoCtx(context.Background(), "DLQ event retried successfully", "id", id)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"status":    "retry_scheduled",
		"id":        id,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}
