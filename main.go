package main 
import(
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
	
	"goeduflow/config"
	"goeduflow/routes"
)
func main(){
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: Error loading .env file (this is normal in production)")
	}
	
	connStr := os.Getenv("SUPABASE_URL")
	if connStr == "" {
		log.Fatal("SUPABASE_URL environment variable is required")
	}
	
	fmt.Println("Supabase Connection:", connStr)
	
	// Connect to database with retry logic
	config.ConnecttoDB()
	
	// Create router
	r := mux.NewRouter()
	routes.RegisterRoutes(r)

	r.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		if err := config.HealthCheck(); err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte(fmt.Sprintf("Database health check failed: %v", err)))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods("GET")

	// Basic endpoint to verify the backend
	r.HandleFunc("/api/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello from Go backend"))
	}).Methods("GET")

	// CORS configuration
	c := cors.New(cors.Options{
		AllowedOrigins: []string{os.Getenv("EduFlow_API")},
		AllowedMethods: []string{"GET","POST","PUT","DELETE"},
		AllowedHeaders: []string{"Authorization","Content-Type"},
		AllowCredentials: true,
	})
	handler := c.Handler(r)
	
	// Create server with timeouts
	server := &http.Server{
		Addr:         ":8080",
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	
	// Graceful shutdown
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan
		
		log.Println("Shutting down server...")
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		
		if err := server.Shutdown(ctx); err != nil {
			log.Printf("Server forced to shutdown: %v", err)
		}
		
		// Close database connection
		if config.GetDB() != nil {
			config.GetDB().Close()
		}
	}()
	
	fmt.Println("Server starting on port 8080...")
	
	// Start server
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal("Server failed to start:", err)
	}
}