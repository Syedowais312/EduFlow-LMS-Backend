package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"

)

var DB *sql.DB

func ConnecttoDB() {
	
	var err error
	DB, err = sql.Open("postgres", os.Getenv("SUPABASE_URL"))
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
