package handlers

import (
	"database/sql"
	"encoding/json"
	"godesk/database"
	"godesk/internal"
	"log"
	"net/http"
	"strings"
)

type Signupdata struct {
	Name     string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type ApiResponse struct {
	Message  string `json:"message"`
	UserName string `json:"username"`
}

func SignUpUser(w http.ResponseWriter, r *http.Request) {
	// decode the signup data form and put it in struct
	var userdata Signupdata
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&userdata); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}
	if userdata.Name == "" || userdata.Email == "" || userdata.Password == "" || userdata.Role == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}
	// generate the password by bcrypt.
	encryptedPassword, err := internal.GenerateTheHashPassword(userdata.Password)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}
	query := `INSERT INTO users (name, email, password, role) VALUES (?, ?, ?, ?)`
	_, err = database.DB.Exec(query, userdata.Name, userdata.Email, encryptedPassword, userdata.Role)
	if err != nil {
		http.Error(w, "Failed to insert user", http.StatusInternalServerError)
		return
	}
	resp := ApiResponse{Message: "User created successfully", UserName: userdata.Name}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
	log.Printf("New user signed up: %s with role %s", userdata.Name, userdata.Role)
}

type Logindata struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type LoginApiResponse struct {
	Message string `json:"message"`
	Token   string `json:"token"`
}

func LoginUser(w http.ResponseWriter, r *http.Request) {
	var logindata Logindata
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&logindata); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}
	if logindata.Email == "" || logindata.Password == "" {
		http.Error(w, "Email and password are required", http.StatusBadRequest)
		return
	}
	query := `SELECT id, name, password, role FROM users WHERE email = ?`
	row := database.DB.QueryRow(query, logindata.Email)
	var userid int
	var username, hashedPassword, role string
	err := row.Scan(&userid, &username, &hashedPassword, &role)
	if err == sql.ErrNoRows {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	if !internal.ValidateThePassword(hashedPassword, logindata.Password) {
		http.Error(w, "Incorrect password", http.StatusUnauthorized)
		return
	}
	token, err := internal.CreateSession(userid, role)
	if err != nil {
		http.Error(w, "Session creation failed", http.StatusInternalServerError)
		return
	}
	response := LoginApiResponse{Message: "Login successful", Token: token}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
	log.Printf("User %s (Role: %s) logged in successfully", username, role)
}

// for logout since the user is loggedIn we get the header and then delete the session
func LogoutUser(w http.ResponseWriter, r *http.Request) {
	header := r.Header.Get("Authorization")
	if header == "" {
		http.Error(w, "Authorization header required", http.StatusUnauthorized)
		return
	}
	tokenParts := strings.Split(header, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		http.Error(w, "Invalid Authorization format", http.StatusUnauthorized)
		return
	}
	token := tokenParts[1]
	internal.DeleteSession(token)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Logged out successfully"))
}
