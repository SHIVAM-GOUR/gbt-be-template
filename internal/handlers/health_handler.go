package handlers

import (
	"net/http"
	"time"

	"gbt-be-template/internal/repository"
	"gbt-be-template/pkg/logger"
	"gbt-be-template/pkg/utils"
)

// HealthHandler handles health check requests
type HealthHandler struct {
	db  *repository.Database
	log *logger.Logger
}

// NewHealthHandler creates a new health handler
func NewHealthHandler(db *repository.Database, log *logger.Logger) *HealthHandler {
	return &HealthHandler{
		db:  db,
		log: log,
	}
}

// Health handles GET /health
func (h *HealthHandler) Health(w http.ResponseWriter, r *http.Request) {
	status := "healthy"
	statusCode := http.StatusOK
	
	// Check database health
	dbStatus := "healthy"
	if err := h.db.Health(); err != nil {
		h.log.WithError(err).Error("Database health check failed")
		dbStatus = "unhealthy"
		status = "unhealthy"
		statusCode = http.StatusServiceUnavailable
	}

	// Get database stats
	dbStats := h.db.GetStats()

	response := map[string]interface{}{
		"status":    status,
		"timestamp": time.Now().UTC(),
		"version":   "1.0.0",
		"services": map[string]interface{}{
			"database": map[string]interface{}{
				"status": dbStatus,
				"stats":  dbStats,
			},
		},
	}

	if status == "healthy" {
		utils.WriteSuccessResponse(w, statusCode, "Service is healthy", response)
	} else {
		utils.WriteErrorResponse(w, statusCode, "Service is unhealthy", response)
	}
}

// Ready handles GET /ready
func (h *HealthHandler) Ready(w http.ResponseWriter, r *http.Request) {
	// Check if all dependencies are ready
	ready := true
	
	// Check database connection
	if err := h.db.Health(); err != nil {
		h.log.WithError(err).Error("Database readiness check failed")
		ready = false
	}

	if ready {
		utils.WriteSuccessResponse(w, http.StatusOK, "Service is ready", map[string]interface{}{
			"ready":     true,
			"timestamp": time.Now().UTC(),
		})
	} else {
		utils.WriteErrorResponse(w, http.StatusServiceUnavailable, "Service is not ready", map[string]interface{}{
			"ready":     false,
			"timestamp": time.Now().UTC(),
		})
	}
}

// Live handles GET /live
func (h *HealthHandler) Live(w http.ResponseWriter, r *http.Request) {
	// Simple liveness check - if we can respond, we're alive
	utils.WriteSuccessResponse(w, http.StatusOK, "Service is alive", map[string]interface{}{
		"alive":     true,
		"timestamp": time.Now().UTC(),
	})
}
