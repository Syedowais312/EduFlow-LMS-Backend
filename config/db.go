package config

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
	"goeduflow/utils"
)

var DB *sql.DB

func ConnecttoDB() {
	connStr := os.Getenv("SUPABASE_URL")
	if connStr == "" {
		log.Fatal("SUPABASE_URL environment variable is not set")
	}

	// Validate connection string
	if err := utils.ValidateConnectionString(connStr); err != nil {
		log.Fatal("Invalid connection string:", err)
	}

	// Optimize connection string for cloud deployment
	connStr = utils.OptimizeConnectionString(connStr)
	
	// Log masked connection string for debugging
	log.Printf("Connecting to database: %s", utils.MaskConnectionString(connStr))

	var err error
	maxRetries := 5
	retryDelay := time.Second * 2

	for i := 0; i < maxRetries; i++ {
		DB, err = sql.Open("postgres", connStr)
		if err != nil {
			log.Printf("Attempt %d: Failed to open database connection: %v", i+1, err)
			if i < maxRetries-1 {
				time.Sleep(retryDelay)
				retryDelay *= 2 // Exponential backoff
				continue
			}
			log.Fatal("Failed to open database connection after all retries:", err)
		}

		// Configure connection pool for cloud deployment
		DB.SetMaxOpenConns(25)                 // Maximum number of open connections
		DB.SetMaxIdleConns(5)                  // Maximum number of idle connections
		DB.SetConnMaxLifetime(time.Hour)       // Maximum lifetime of a connection
		DB.SetConnMaxIdleTime(time.Minute * 5) // Maximum idle time of a connection

		// Test the connection with timeout
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		
		err = DB.PingContext(ctx)
		if err != nil {
			log.Printf("Attempt %d: Failed to ping database: %v", i+1, err)
			DB.Close()
			if i < maxRetries-1 {
				time.Sleep(retryDelay)
				retryDelay *= 2
				continue
			}
			log.Fatal("Failed to connect to database after all retries:", err)
		}

		fmt.Println("Connected to Database successfully")
		return
	}
}

// HealthCheck verifies database connectivity
func HealthCheck() error {
	if DB == nil {
		return fmt.Errorf("database connection is nil")
	}
	
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	return DB.PingContext(ctx)
}

// GetDB returns the database connection
func GetDB() *sql.DB {
	return DB
}
