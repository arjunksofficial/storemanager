package handler

import (
	"fmt"
	"net/http"
	"time"

	"storemanager/internal/common/responsehelper"

	"github.com/go-chi/chi"
)

// HealthHandler is handler object for health check
type HealthHandler struct {
	StartTime *time.Time
}

// NewHealth gives a new HealthHandler
func NewHealth(startTime *time.Time) *HealthHandler {
	return &HealthHandler{
		StartTime: startTime,
	}
}

// SetHealthRoutes creates routes for health check
func SetHealthRoutes(r chi.Router, hH *HealthHandler) {
	r.MethodFunc(http.MethodGet, "/", hH.GetHealthStatus)
}

// GetHealthStatus gives current server status
func (hH *HealthHandler) GetHealthStatus(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("server is up for %v, ie Up Since %v", time.Since(*hH.StartTime), *hH.StartTime)
	responsehelper.RespondAsJSON(http.StatusOK, w, responsehelper.CommonResponse{
		Message: message,
		Status:  "success",
	})
}
