package models 
import "time"
type Assignment struct{
	ID int `json:"id"`
	Title string `json:"title"`
	Description string `json:"description"`
	DueDate time.Time `json:"due_date"`
	Teacher_ID int `json:"teacher_id"`
}