package controllers

import (
	"encoding/json"
	"goeduflow/config"
	"goeduflow/models"
	"net/http"
	"time"
	
)
func SubmissionAssignment(w http.ResponseWriter,r *http.Request){
	var submission models.Submission
	err := json.NewDecoder(r.Body).Decode(&submission)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	submission.SubmittedAt=time.Now().Truncate(time.Second)
	query:=`INSERT INTO submissions (file_path,content,assignment_id,student_id,submitted_at)
	            VALUES ($1, $2, $3, $4,$5) RETURNING id`
    err=config.DB.QueryRow(query,submission.File_path,submission.Content,submission.AssignmentID,submission.StudentID,submission.SubmittedAt).Scan(&submission.ID)
if err!=nil{
	http.Error(w, err.Error(),http.StatusInternalServerError)
	return
}
w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
json.NewEncoder(w).Encode(submission)

}
