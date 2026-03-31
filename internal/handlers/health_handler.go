package handlers

import (
	"net/http"

	"go-ecommerce/internal/database"
	apierrors "go-ecommerce/internal/errors"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// Health godoc
// @Summary Health check
// @Description Check if the API is running
// @Tags health
// @Produce json
// @Success 200 {object} map[string]string
// @Router /health [get]
func (h *HealthHandler) Health(w http.ResponseWriter, r *http.Request) {
	apierrors.RespondWithJSON(w, http.StatusOK, map[string]string{
		"status": "ok",
	})
}

// Ready godoc
// @Summary Readiness check
// @Description Check if the API and database are ready
// @Tags health
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 503 {object} map[string]string
// @Router /ready [get]
func (h *HealthHandler) Ready(w http.ResponseWriter, r *http.Request) {
	if database.DB == nil {
		apierrors.RespondWithError(w, http.StatusServiceUnavailable, apierrors.New(apierrors.ErrCodeInternalError, "Database not connected"))
		return
	}

	if err := database.DB.Ping(); err != nil {
		apierrors.RespondWithError(w, http.StatusServiceUnavailable, apierrors.New(apierrors.ErrCodeDatabaseError, "Database not ready"))
		return
	}

	apierrors.RespondWithJSON(w, http.StatusOK, map[string]string{
		"status":   "ready",
		"database": "connected",
	})
}
