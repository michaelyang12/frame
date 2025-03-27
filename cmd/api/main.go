package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/michaelyang12/frame/internal/handlers"
	"github.com/michaelyang12/frame/internal/middleware"
	"github.com/rs/cors"
)

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*") // Allow all origins (more restrictive: specific domain)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight OPTIONS requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Pass to next handler
		next.ServeHTTP(w, r)
	})
}

func main() {
	port := 5050
	// Set up server with timeouts
	server := &http.Server{
		Addr:         fmt.Sprintf("0.0.0.0:%s", port),
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      setupRoutes(),
	}

	// Start server
	log.Println("Image Processing API running on port %s", port)
	log.Fatal(server.ListenAndServe())
}

func setupRoutes() http.Handler {
	// Define routes
	mux := http.NewServeMux()
	mux.HandleFunc("/", handlers.HandleRoot)
	mux.HandleFunc("/resize", handlers.HandleResize)
	mux.HandleFunc("/convert", handlers.HandleConvert)
	mux.HandleFunc("/trim", handlers.HandleTrim)
	mux.HandleFunc("/health", handlers.HandleHealth)

	// Apply middleware
	handler := middleware.Logging(mux)

	// Apply CORS middleware
	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	return corsMiddleware.Handler(handler)
}
