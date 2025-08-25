package models
 type User struct{
	ID int `json:"id"`
	Firstname string `json:"firstname"`
	Lastname string `json:"lastname"`
	Email string `json:"email"`
	Grade int `json:"grade"`
	School string `json:"school"`
	Password string `json:"password"`
	Role string `json:"role"`
	CreatedAt string `json:"created_at"`
	
 }