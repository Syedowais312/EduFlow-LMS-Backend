package main 
import(
	"fmt"
	"log"
	"os"
	"net/http"
	"github.com/gorilla/mux"
	"goeduflow/config"
	"goeduflow/routes"
	"github.com/rs/cors"
	"github.com/joho/godotenv"
)
func main(){
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	connStr := os.Getenv("SUPABASE_URL")
	fmt.Println("Supabase Connection:", connStr)
	// Connect to database
	config.ConnecttoDB()
	
	// Create router
	r := mux.NewRouter()
	routes.RegisterRoutes(r)

	//just to verify the backend
	r.HandleFunc("/api/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello from Go backend "))
	}).Methods("GET")

	c:=cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:3000"},
		AllowedMethods: []string{"GET","POST","PUT","DELETE"},
		AllowedHeaders: []string{"Authorization","Content-Type"},
		AllowCredentials: true,
	})
	handler:=c.Handler(r)
	
	fmt.Println("Server starting on port 8080...")
	
	// Start server with error handling
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}