package controllers

import (
	"encoding/json"
	"net/http"
	"time"

	"goeduflow/config"
	"goeduflow/middleware"
	"goeduflow/models"
)

// CreateAssignment handler
func CreateAssignment(w http.ResponseWriter, r *http.Request) {
	var assignment models.Assignment
	err := json.NewDecoder(r.Body).Decode(&assignment)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	query := `INSERT INTO assignments (title, description, due_date, teacher_id, school, subject) 
	          VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
	err = config.DB.QueryRow(
		query,
		assignment.Title,
		assignment.Description,
		assignment.DueDate,
		assignment.Teacher_ID,
		assignment.School,
		assignment.Subject,
	).Scan(&assignment.ID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	assignment.DueDate = assignment.DueDate.Local().Truncate(time.Second)
	json.NewEncoder(w).Encode(assignment)
}

// FetchAssignments handler
func FetchAssignments(w http.ResponseWriter, r *http.Request) {
	// Get school name from middleware context
	school, ok := r.Context().Value(authmiddleware.SchoolKey).(string)
	if !ok || school == "" {
		jsonError(w, http.StatusUnauthorized, "School not found in token")
		return
	}

	rows, err := config.DB.Query(`
		SELECT id, title, description, due_date, subject, school
		FROM assignments 
		WHERE school=$1`,
		school,
	)
	if err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer rows.Close()

	var assignments []models.Assignment
	for rows.Next() {
		var a models.Assignment
		if err := rows.Scan(&a.ID, &a.Title, &a.Description, &a.DueDate, &a.Subject, &a.School); err != nil {
			jsonError(w, http.StatusInternalServerError, err.Error())
			return
		}
		assignments = append(assignments, a)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(assignments)
}
func jsonError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{
		"error": message,
	})
}
