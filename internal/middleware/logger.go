package middleware

import (
	"log"
	"net/http"
	"time"
)

// Logging middleware logs information about each request
func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Create a custom response writer to capture status code
		rw := &responseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		// Process request
		next.ServeHTTP(rw, r)

		// Log after request is processed
		duration := time.Since(start)
		log.Printf(
			"[%s] %s %s %d %s",
			r.Method,
			r.RequestURI,
			r.RemoteAddr,
			rw.statusCode,
			duration,
		)
	})
}

// Recovery middleware recovers from panics and logs the error
// func Recovery(next http.Handler) http.Handler {
// }

// Timeout middleware adds a timeout to the request context
// func Timeout(duration time.Duration) func(http.Handler) http.Handler {
// }

// // CORS middleware adds Cross-Origin Resource Sharing headers
// func CORS(next http.Handler) http.Handler {

// }

// // RateLimit middleware implements a simple rate limiting strategy
// // Note: For production, consider using a more sophisticated rate limiter
// func RateLimit(requestsPerMinute int) func(http.Handler) http.Handler {
// 	// Simple in-memory store of IP addresses and request timestamps

// }

// // Security middleware adds basic security headers
// func Security(next http.Handler) http.Handler {

// }

// // responseWriter is a custom ResponseWriter that captures the status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

// // WriteHeader captures the status code before writing it
// func (rw *responseWriter) WriteHeader(code int) {
// 	rw.statusCode = code
// 	rw.ResponseWriter.WriteHeader(code)
// }
