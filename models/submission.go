package models
import "time"
type Submission struct{
	ID int `json:"id"`
	AssignmentID int `json:"assignment_id"`
	StudentID int `json:"student_id"`
     SubmittedAt time.Time `json:"submitted_at"`
	 Content string `json:"content"`
	 File_path string `json:"file_path"`
}