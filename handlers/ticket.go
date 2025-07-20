package handlers

import (
	"database/sql"
	"encoding/json"
	"godesk/appcontext"
	"godesk/database"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type TicketResponse struct {
	ID          int            `json:"id"`
	Token       int            `json:"ticket_token"`
	UserID      int            `json:"user_id"`
	CategoryID  int            `json:"category_id"`
	Title       string         `json:"title"`
	Description string         `json:"description"`
	Status      string         `json:"status"`
	Priority    sql.NullString `json:"priority"`
	AssignedTo  sql.NullString `json:"assigned_to"`
	CreatedAt   string         `json:"created_at"`
	UpdatedAt   string         `json:"updated_at"`
}

func GetTickets(w http.ResponseWriter, r *http.Request) {
	// use the  same method for admins and employee
	// getting the user and role by context via the middleware
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

	// making the global variable for query and error
	var rows *sql.Rows
	var err error

	if role == "admin" {
		query := "SELECT id, ticket_token, user_id, category_id, title, description, status, priority, assigned_to, created_at, updated_at FROM tickets ORDER BY created_at DESC;"
		rows, err = database.DB.Query(query)
	} else {
		query := "SELECT id, ticket_token, user_id, category_id, title, description, status, priority, assigned_to, created_at, updated_at FROM tickets WHERE user_id = ? ORDER BY created_at DESC;"
		rows, err = database.DB.Query(query, userID)
	}
	if err != nil {
		http.Error(w, "tickets cannot be fetched from database error...", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var tickets []TicketResponse
	for rows.Next() {
		var t TicketResponse
		err := rows.Scan(&t.ID, &t.Token, &t.UserID, &t.CategoryID, &t.Title, &t.Description, &t.Status, &t.Priority, &t.AssignedTo, &t.CreatedAt, &t.UpdatedAt)
		if err != nil {
			log.Printf("Error scanning ticket row: %v", err)
			continue
		}
		tickets = append(tickets, t)
	}
	if err = rows.Err(); err != nil {
		http.Error(w, "Error iterating ticket rows", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&tickets)
}

// create ticket (post) function.
func CreateTicket(w http.ResponseWriter, r *http.Request) {
	var ticketData struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		CategoryID  int    `json:"category_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&ticketData); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}
	userID, ok := r.Context().Value(appcontext.UserIDKey).(int)
	if !ok {
		http.Error(w, "Could not retrieve user ID from context", http.StatusInternalServerError)
		return
	}
	tokenID := rand.Intn(90000) + 10000

	query := "INSERT INTO tickets(ticket_token, user_id, category_id, title, description) VALUES(?,?,?,?,?)"
	_, err := database.DB.Exec(query, tokenID, userID, ticketData.CategoryID, ticketData.Title, ticketData.Description)
	if err != nil {
		http.Error(w, "Failed to create ticket", http.StatusInternalServerError)
		log.Println("Error creating ticket: ", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Ticket created successfully"})
}

// delete ticket function.
func DeleteTicket(w http.ResponseWriter, r *http.Request) {
	// get the role and user id
	role, ok := r.Context().Value(appcontext.UserRoleKey).(string)
	if !ok {
		http.Error(w, "Could not retrieve user role from context", http.StatusInternalServerError)
		return
	}
	userID, ok := r.Context().Value(appcontext.UserIDKey).(int)
	if !ok {
		http.Error(w, "Could not retrieve user ID from context", http.StatusInternalServerError)
		return
	}
	vars := mux.Vars(r)
	ticketIDStr := vars["id"]
	ticketID, err := strconv.Atoi(ticketIDStr)
	if err != nil {
		http.Error(w, "Invalid ticket ID", http.StatusBadRequest)
		return
	}

	// var res sql.Result
	if role == "admin" {
		_, err = database.DB.Exec(`DELETE FROM tickets WHERE id = ?`, ticketID)
	} else if role == "employee" {
		_, err = database.DB.Exec(`DELETE FROM tickets WHERE id = ? AND user_id = ?`, ticketID, userID)
	} else {
		http.Error(w, "Error deleting ticket", http.StatusInternalServerError)
		return
	}

	// if err != nil {
	// 	http.Error(w, "Error deleting ticket", http.StatusInternalServerError)
	// 	return
	// }
	// rowsAffected, _ := res.RowsAffected()
	// if rowsAffected == 0 {
	// 	http.Error(w, "Ticket not found or you do not have permission to delete it", http.StatusNotFound)
	// 	return
	// }
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Ticket deleted successfully"))
}

// function for assign ticket.
func AssignTicket(w http.ResponseWriter, r *http.Request) {
	role, ok := r.Context().Value(appcontext.UserRoleKey).(string)
	if !ok {
		http.Error(w, "Could not retrieve user role from context", http.StatusInternalServerError)
		return
	}
	if role != "admin" {
		http.Error(w, "Forbidden: Only admins can assign tickets", http.StatusForbidden)
		return
	}

	var data struct {
		TicketID   int    `json:"ticket_id"`
		AssignedTo string `json:"assigned_to"`
	}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	query := `UPDATE tickets SET assigned_to = ? WHERE id = ?`
	_, err := database.DB.Exec(query, data.AssignedTo, data.TicketID)
	if err != nil {
		http.Error(w, "Error assigning ticket", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Ticket assigned successfully"))
}

func SetPriority(w http.ResponseWriter, r *http.Request) {
	role, ok := r.Context().Value(appcontext.UserRoleKey).(string)
	if !ok {
		http.Error(w, "Could not retrieve user role from context", http.StatusInternalServerError)
		return
	}
	if role != "admin" {
		http.Error(w, "Forbidden: Only admins can assign tickets", http.StatusForbidden)
		return
	}

	var data struct {
		TicketID int    `json:"ticket_id"`
		Priority string `json:"priority"`
	}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	query := `UPDATE tickets SET priority = ? WHERE id = ?`
	_, err := database.DB.Exec(query, data.Priority, data.TicketID)
	if err != nil {
		http.Error(w, "Error setting priority", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Priority updated successfully"))
}

func SetStatus(w http.ResponseWriter, r *http.Request) {
	role, ok := r.Context().Value(appcontext.UserRoleKey).(string)
	if !ok {
		http.Error(w, "Could not retrieve user role from context", http.StatusInternalServerError)
		return
	}
	if role != "admin" {
		http.Error(w, "Forbidden: Only admins can assign tickets", http.StatusForbidden)
		return
	}

	var data struct {
		TicketID int    `json:"ticket_id"`
		Status   string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	query := `UPDATE tickets SET status = ? WHERE id = ?`
	_, err := database.DB.Exec(query, data.Status, data.TicketID)
	if err != nil {
		http.Error(w, "Error setting status", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Status updated successfully"))
}
