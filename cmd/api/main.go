package main

import (
	"log"
	"net/http"
	"time"

	"github.com/michaelyang12/frame/internal/handlers"
	"github.com/michaelyang12/frame/internal/middleware"
)

func main() {
	// Set up server with timeouts
	server := &http.Server{
		Addr:         "0.0.0.0:8080",
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      setupRoutes(),
	}

	// Start server
	log.Println("Image Processing API running on port 8080")
	log.Fatal(server.ListenAndServe())
}

func setupRoutes() http.Handler {
	// Define routes
	mux := http.NewServeMux()
	mux.HandleFunc("/", handlers.HandleRoot)
	mux.HandleFunc("/resize", handlers.HandleResize)
	mux.HandleFunc("/convert", handlers.HandleConvert)
	mux.HandleFunc("/health", handlers.HandleHealth)

	// Apply middleware
	return middleware.Logging(mux)
}
