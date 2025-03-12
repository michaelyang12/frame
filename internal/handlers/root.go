package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/michaelyang12/frame/internal/models"
)

// HandleRoot returns API status
func HandleRoot(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	resp := models.Response{
		Success: true,
		Message: "Image Processing API is running",
		Data: map[string]string{
			"version":   "1.0.0",
			"endpoints": "/resize, /convert, /health",
		},
	}

	sendJSONResponse(w, resp, http.StatusOK)
}

// sendJSONResponse standardizes JSON responses
func sendJSONResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Error encoding JSON: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
