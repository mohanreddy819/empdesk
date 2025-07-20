package handlers

import (
	"encoding/json"
	"godesk/appcontext"
	"godesk/database"
	"net/http"
)

func ActivityLog(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(appcontext.UserIDKey).(int)
	if !ok {
		http.Error(w, "Could not retrieve user ID from context", http.StatusInternalServerError)
		return
	}
	role, ok := r.Context().Value(appcontext.UserRoleKey).(string)
	if !ok {
		http.Error(w, "Could not retrieve user role from context", http.StatusInternalServerError)
		return
	}

	if role != "admin" {
		http.Error(w, "Forbidden: Only admins can log activity", http.StatusForbidden)
		return
	}

	var input struct {
		Action string `json:"action"`
	}
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil || input.Action == "" {
		http.Error(w, "Invalid input. 'action' is required.", http.StatusBadRequest)
		return
	}

	_, err = database.DB.Exec(`INSERT INTO activity_logs (user_id, action) VALUES (?, ?)`, userID, input.Action)
	if err != nil {
		http.Error(w, "Failed to log activity", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Activity logged successfully"))
}
