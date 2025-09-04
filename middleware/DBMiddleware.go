package middleware

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"goeduflow/config"
)

// DBMiddleware handles database connection issues and retries
func DBMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check database health before processing request
		if err := config.HealthCheck(); err != nil {
			// Log the error but don't fail the request immediately
			fmt.Printf("Database health check failed: %v\n", err)
			
			// Try to reconnect once
			config.ConnecttoDB()
			
			// Check again after reconnection attempt
			if err := config.HealthCheck(); err != nil {
				http.Error(w, "Database temporarily unavailable", http.StatusServiceUnavailable)
				return
			}
		}
		
		// Add database connection to request context
		ctx := context.WithValue(r.Context(), "db", config.GetDB())
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetDBFromContext retrieves database connection from request context
func GetDBFromContext(r *http.Request) *sql.DB {
	if db, ok := r.Context().Value("db").(*sql.DB); ok {
		return db
	}
	return config.GetDB()
}

// WithTimeout creates a context with timeout for database operations
func WithTimeout(r *http.Request, timeout time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(r.Context(), timeout)
}
