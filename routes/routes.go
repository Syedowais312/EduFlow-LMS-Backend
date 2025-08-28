package routes

import (
	"goeduflow/controllers"
	"github.com/gorilla/mux"
)
func RegisterRoutes(r *mux.Router){
	r.HandleFunc("/signup",controllers.Signup).Methods("POST")
r.HandleFunc("/login", controllers.Login).Methods("POST")
r.HandleFunc("/assignments",controllers.CreateAssignment).Methods("POST")
r.HandleFunc("/submission",controllers.SubmissionAssignment).Methods("POST")
r.HandleFunc("/fetchAssignment",controllers.FetchAssignments).Methods("GET")
}