package utils

import (
	"fmt"
	"net/url"
	"os"
	"strings"
)

// OptimizeConnectionString adds optimal parameters for cloud deployment
func OptimizeConnectionString(connStr string) string {
	// Parse the connection string
	u, err := url.Parse(connStr)
	if err != nil {
		return connStr // Return original if parsing fails
	}

	// Get existing query parameters
	query := u.Query()

	// Add optimized parameters for cloud deployment
	query.Set("sslmode", "require")
	query.Set("connect_timeout", "10")        // 10 second connection timeout
	query.Set("statement_timeout", "30000")   // 30 second statement timeout
	query.Set("idle_in_transaction_session_timeout", "30000") // 30 second idle timeout
	query.Set("tcp_keepalives_idle", "600")   // 10 minutes
	query.Set("tcp_keepalives_interval", "30") // 30 seconds
	query.Set("tcp_keepalives_count", "3")    // 3 keepalive probes

	// Rebuild the URL
	u.RawQuery = query.Encode()
	return u.String()
}

// GetOptimizedConnectionString returns the optimized connection string from environment
func GetOptimizedConnectionString() string {
	connStr := os.Getenv("SUPABASE_URL")
	if connStr == "" {
		return ""
	}
	return OptimizeConnectionString(connStr)
}

// ValidateConnectionString checks if the connection string is valid
func ValidateConnectionString(connStr string) error {
	if connStr == "" {
		return fmt.Errorf("connection string is empty")
	}

	u, err := url.Parse(connStr)
	if err != nil {
		return fmt.Errorf("invalid connection string format: %v", err)
	}

	if u.Scheme != "postgresql" && u.Scheme != "postgres" {
		return fmt.Errorf("unsupported database scheme: %s", u.Scheme)
	}

	if u.Host == "" {
		return fmt.Errorf("missing host in connection string")
	}

	return nil
}

// MaskConnectionString masks sensitive information in connection string for logging
func MaskConnectionString(connStr string) string {
	u, err := url.Parse(connStr)
	if err != nil {
		return "invalid_connection_string"
	}

	// Mask password if present
	if u.User != nil {
		u.User = url.User(u.User.Username())
	}

	// Mask host if it contains sensitive information
	if strings.Contains(u.Host, "supabase") {
		parts := strings.Split(u.Host, ".")
		if len(parts) > 2 {
			parts[1] = "***"
			u.Host = strings.Join(parts, ".")
		}
	}

	return u.String()
}
