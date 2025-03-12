package handlers

import (
	"net/http"
	"time"

	"github.com/michaelyang12/frame/internal/models"
)

func HandleHealth(w http.ResponseWriter, r *http.Request) {
	resp := models.Response{
		Success: true,
		Message: "Service is healthy",
		Data: map[string]interface{}{
			"uptime": time.Now().Format(time.RFC3339),
			"status": "operational",
		},
	}

	sendJSONResponse(w, resp, http.StatusOK)
}
