package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"gbt-be-template/pkg/logger"
	"gbt-be-template/pkg/utils"
)

// Recovery middleware recovers from panics and logs them
func Recovery(log *logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					// Log the panic with stack trace
					log.WithFields(map[string]interface{}{
						"error":      fmt.Sprintf("%v", err),
						"stack":      string(debug.Stack()),
						"method":     r.Method,
						"path":       r.URL.Path,
						"user_agent": r.UserAgent(),
						"ip":         getClientIP(r),
						"type":       "panic",
					}).Error("Panic recovered")

					// Return 500 Internal Server Error
					utils.WriteErrorResponse(w, http.StatusInternalServerError, "Internal server error", nil)
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
