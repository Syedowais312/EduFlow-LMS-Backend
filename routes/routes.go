package routes

import (
	"goeduflow/controllers"
	"goeduflow/middleware"
	"github.com/gorilla/mux"
)

func RegisterRoutes(r *mux.Router){
	// Apply database middleware to all routes
	r.Use(middleware.DBMiddleware)
	
	// Auth routes
	r.HandleFunc("/signup", controllers.Signup).Methods("POST")
	r.HandleFunc("/login", controllers.Login).Methods("POST")
	
	// Assignment routes
	r.HandleFunc("/assignments", controllers.CreateAssignment).Methods("POST")
	r.HandleFunc("/submission", controllers.SubmissionAssignment).Methods("POST")
	r.HandleFunc("/fetchAssignment", controllers.FetchAssignments).Methods("GET")
}