package controllers

import (
	"encoding/json"
	"goeduflow/config"
	"goeduflow/models"
	"goeduflow/utils"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

// function for signup
func Signup(w http.ResponseWriter, r *http.Request) {
	var user models.User
	json.NewDecoder(r.Body).Decode(&user)

	// validate input
	if user.Email == "" || user.Password == "" {
		jsonResponse(w, http.StatusBadRequest, "Email and password required")
		return
	}

	// hash password
	hashedPasswords, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, "Error hashing password")
		return
	}

	// insert into DB
	err = config.DB.QueryRow(
    "INSERT INTO users (firstname, lastname, email, password, role, grade, school) VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING id",
    user.Firstname, user.Lastname, user.Email, string(hashedPasswords), user.Role, user.Grade, user.School,
).Scan(&user.ID)

	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			jsonResponse(w, http.StatusBadRequest, "User already exists with this email")
			return
		}
		jsonResponse(w, http.StatusInternalServerError, "Error saving user: "+err.Error())
		return
	}
	token,err:=utils.GenerateJWT(user.Firstname+" "+user.Lastname,user.Email,user.School,user.ID)
	if(err!=nil){
		jsonResponse(w, http.StatusInternalServerError,"Error generating token")
	}

	    response := map[string]interface{}{
        "token": token,
        "user": map[string]interface{}{
			"id":        user.ID,
            "firstname": user.Firstname,
            "lastname":  user.Lastname,
            "email":     user.Email,
            "role":      user.Role,
            "grade":     user.Grade,
            "school":    user.School,
        },
    }

	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// function for login
func Login(w http.ResponseWriter, r *http.Request) {
	var user models.User
	json.NewDecoder(r.Body).Decode(&user)

	var dbUser models.User
err := config.DB.QueryRow(
    "SELECT id, firstname, lastname, email, role, password, grade, school FROM users WHERE email=$1",
    user.Email,
).Scan(&dbUser.ID, &dbUser.Firstname, &dbUser.Lastname, &dbUser.Email, &dbUser.Role, &dbUser.Password, &dbUser.Grade, &dbUser.School)

	if err != nil {
		jsonResponse(w, http.StatusUnauthorized, "Invalid email or user not found")
		return
	}

	// compare password
	if err := bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password)); err != nil {
		jsonResponse(w, http.StatusUnauthorized, "Invalid password")
		return
	}

	// generate JWT
	token, err := utils.GenerateJWT(dbUser.Firstname+" "+dbUser.Lastname,dbUser.Email,dbUser.School,dbUser.ID)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, "Error generating token")
		return
	}

	// response with token + user
	response := map[string]interface{}{
    "token": token,
    "user": map[string]interface{}{
			"id":        dbUser.ID,
        "firstname": dbUser.Firstname,
        "lastname":  dbUser.Lastname,
        "email":     dbUser.Email,
        "role":      dbUser.Role,
        "grade":     dbUser.Grade,
        "school":    dbUser.School,
    },
}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// helper for sending JSON responses
func jsonResponse(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"message": message})
}
