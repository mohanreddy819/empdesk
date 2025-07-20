package handlers

import (
	"encoding/json"
	"fmt"
	"godesk/database"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type CommentInput struct {
	Comments string `json:"comments"`
}

func AddComment(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int)
	ticketIDStr := mux.Vars(r)["id"]
	ticketID, err := strconv.Atoi(ticketIDStr)
	if err != nil {
		http.Error(w, "Invalid ticket ID", http.StatusBadRequest)
		return
	}

	var input CommentInput
	err = json.NewDecoder(r.Body).Decode(&input)
	if err != nil || input.Comments == "" {
		http.Error(w, "Invalid comment input", http.StatusBadRequest)
		return
	}

	_, err = database.DB.Exec(`
		INSERT INTO ticket_comments (ticket_id, user_id, comments)
		VALUES (?, ?, ?)`,
		ticketID, userID, input.Comments,
	)

	if err != nil {
		http.Error(w, "Failed to insert comment", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintln(w, "Comment added successfully")
}
