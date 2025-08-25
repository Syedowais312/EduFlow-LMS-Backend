package controllers

import (
	"encoding/json"
	"net/http"
	"time"

	"goeduflow/config"
	"goeduflow/models"
)

func CreateAssignment(w http.ResponseWriter, r *http.Request) {
	var assignment models.Assignment
	err := json.NewDecoder(r.Body).Decode(&assignment)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	query := `INSERT INTO assignments (title, description, due_date, teacher_id) 
	          VALUES ($1, $2, $3, $4) RETURNING id`
	err = config.DB.QueryRow(query, assignment.Title, assignment.Description, assignment.DueDate, assignment.Teacher_ID).Scan(&assignment.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	assignment.DueDate = assignment.DueDate.Local().Truncate(time.Second)
	json.NewEncoder(w).Encode(assignment)
}
