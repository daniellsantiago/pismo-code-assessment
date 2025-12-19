package handler

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"github.com/nubank/pismo-code-assessment/internal/infrastructure/http/response"
)

type HealthHandler struct {
	db *sql.DB
}

func NewHealthHandler(db *sql.DB) *HealthHandler {
	return &HealthHandler{db: db}
}

type HealthResponse struct {
	Status   string `json:"status"`
	Database string `json:"database"`
}

func (h *HealthHandler) Check(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	dbStatus := "connected"
	status := "healthy"
	httpStatus := http.StatusOK

	if err := h.db.PingContext(ctx); err != nil {
		dbStatus = "disconnected"
		status = "unhealthy"
		httpStatus = http.StatusServiceUnavailable
	}

	response.JSON(w, httpStatus, HealthResponse{
		Status:   status,
		Database: dbStatus,
	})
}
