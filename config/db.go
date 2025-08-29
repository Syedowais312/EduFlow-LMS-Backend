package config

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func ConnecttoDB() {
	connStr := "postgresql://postgres:Owais@786@db.fxqrwbrbumhjlqowglcx.supabase.co:5432/postgres"
	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Failed to open database connection:", err)
	}
	
	// Test the connection
	err = DB.Ping()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	
	fmt.Println("Connected to Database successfully")
}
