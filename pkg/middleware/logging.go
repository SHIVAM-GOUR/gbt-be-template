package middleware

import (
	"net/http"
	"time"

	"gbt-be-template/pkg/logger"

	"github.com/go-chi/chi/v5/middleware"
)

// responseWriter wraps http.ResponseWriter to capture status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// Logging middleware logs HTTP requests
func Logging(log *logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Wrap the response writer to capture status code
			wrapped := &responseWriter{
				ResponseWriter: w,
				statusCode:     http.StatusOK,
			}

			// Get request ID from context (if using chi's RequestID middleware)
			requestID := middleware.GetReqID(r.Context())

			// Process request
			next.ServeHTTP(wrapped, r)

			// Calculate duration
			duration := time.Since(start).Milliseconds()

			// Get client IP
			ip := getClientIP(r)

			// Log the request
			entry := log.HTTP(
				r.Method,
				r.URL.Path,
				r.UserAgent(),
				ip,
				wrapped.statusCode,
				duration,
			)

			if requestID != "" {
				entry = entry.WithField("request_id", requestID)
			}

			// Log with appropriate level based on status code
			if wrapped.statusCode >= 500 {
				entry.Error("HTTP request completed with server error")
			} else if wrapped.statusCode >= 400 {
				entry.Warn("HTTP request completed with client error")
			} else {
				entry.Info("HTTP request completed")
			}
		})
	}
}

// getClientIP extracts the client IP from the request
func getClientIP(r *http.Request) string {
	// Check X-Forwarded-For header
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		// Take the first IP in the list
		if idx := len(xff); idx > 0 {
			if commaIdx := 0; commaIdx < idx {
				for i, char := range xff {
					if char == ',' {
						commaIdx = i
						break
					}
				}
				if commaIdx > 0 {
					return xff[:commaIdx]
				}
			}
			return xff
		}
	}

	// Check X-Real-IP header
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}

	// Fall back to RemoteAddr
	return r.RemoteAddr
}
