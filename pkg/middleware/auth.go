package middleware

import (
	"context"
	"net/http"
	"strings"

	"gbt-be-template/pkg/logger"
	"gbt-be-template/pkg/utils"
)

// ContextKey is a custom type for context keys to avoid collisions
type ContextKey string

const (
	// UserIDKey is the context key for user ID
	UserIDKey ContextKey = "user_id"
	// UserEmailKey is the context key for user email
	UserEmailKey ContextKey = "user_email"
	// IsAdminKey is the context key for admin status
	IsAdminKey ContextKey = "is_admin"
)

// JWTAuth middleware validates JWT tokens
func JWTAuth(log *logger.Logger, jwtSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get token from Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				log.WithField("path", r.URL.Path).Warn("Missing authorization header")
				utils.WriteErrorResponse(w, http.StatusUnauthorized, "Authorization header required", nil)
				return
			}

			// Check if header starts with "Bearer "
			if !strings.HasPrefix(authHeader, "Bearer ") {
				log.WithField("path", r.URL.Path).Warn("Invalid authorization header format")
				utils.WriteErrorResponse(w, http.StatusUnauthorized, "Invalid authorization header format", nil)
				return
			}

			// Extract token
			token := strings.TrimPrefix(authHeader, "Bearer ")
			if token == "" {
				log.WithField("path", r.URL.Path).Warn("Empty token")
				utils.WriteErrorResponse(w, http.StatusUnauthorized, "Token required", nil)
				return
			}

			// Validate token and extract claims
			claims, err := utils.ValidateJWT(token, jwtSecret)
			if err != nil {
				log.WithError(err).WithField("path", r.URL.Path).Warn("Invalid token")
				utils.WriteErrorResponse(w, http.StatusUnauthorized, "Invalid token", nil)
				return
			}

			// Add user information to context
			ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
			ctx = context.WithValue(ctx, UserEmailKey, claims.Email)
			ctx = context.WithValue(ctx, IsAdminKey, claims.IsAdmin)

			// Continue with the request
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RequireAdmin middleware ensures the user is an admin
func RequireAdmin(log *logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Check if user is admin
			isAdmin, ok := r.Context().Value(IsAdminKey).(bool)
			if !ok || !isAdmin {
				userID := r.Context().Value(UserIDKey)
				log.WithFields(map[string]interface{}{
					"user_id": userID,
					"path":    r.URL.Path,
				}).Warn("Admin access required")
				utils.WriteErrorResponse(w, http.StatusForbidden, "Admin access required", nil)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// OptionalAuth middleware validates JWT tokens but doesn't require them
func OptionalAuth(log *logger.Logger, jwtSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get token from Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				// No token provided, continue without authentication
				next.ServeHTTP(w, r)
				return
			}

			// Check if header starts with "Bearer "
			if !strings.HasPrefix(authHeader, "Bearer ") {
				// Invalid format, continue without authentication
				next.ServeHTTP(w, r)
				return
			}

			// Extract token
			token := strings.TrimPrefix(authHeader, "Bearer ")
			if token == "" {
				// Empty token, continue without authentication
				next.ServeHTTP(w, r)
				return
			}

			// Validate token and extract claims
			claims, err := utils.ValidateJWT(token, jwtSecret)
			if err != nil {
				// Invalid token, continue without authentication
				log.WithError(err).WithField("path", r.URL.Path).Debug("Invalid optional token")
				next.ServeHTTP(w, r)
				return
			}

			// Add user information to context
			ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
			ctx = context.WithValue(ctx, UserEmailKey, claims.Email)
			ctx = context.WithValue(ctx, IsAdminKey, claims.IsAdmin)

			// Continue with the request
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetUserIDFromContext extracts user ID from context
func GetUserIDFromContext(ctx context.Context) (uint, bool) {
	userID, ok := ctx.Value(UserIDKey).(uint)
	return userID, ok
}

// GetUserEmailFromContext extracts user email from context
func GetUserEmailFromContext(ctx context.Context) (string, bool) {
	email, ok := ctx.Value(UserEmailKey).(string)
	return email, ok
}

// GetIsAdminFromContext extracts admin status from context
func GetIsAdminFromContext(ctx context.Context) (bool, bool) {
	isAdmin, ok := ctx.Value(IsAdminKey).(bool)
	return isAdmin, ok
}
