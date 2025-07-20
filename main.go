package main

import (
	"context"
	"fmt"
	"godesk/appcontext"
	"godesk/database"
	"godesk/handlers"
	"godesk/internal"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

// setting up the middleware by taking in the header(authorization) from the request
func SessionAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		// if header == "" {
		// 	http.Error(w, "Authorization header required", http.StatusUnauthorized)
		// 	return
		// }

		// getting the header
		tokenParts := strings.Split(header, " ")
		if len(tokenParts) != 2 {
			http.Error(w, "Invalid Authorization format...", http.StatusUnauthorized)
			return
		} else if tokenParts[0] != "Bearer" {
			http.Error(w, "Invalid Bearer not in Header...", http.StatusUnauthorized)
			return
		}

		token := tokenParts[1]
		sessionData, found := internal.GetSessionData(token)
		if !found {
			http.Error(w, "Invalid session token", http.StatusUnauthorized)
			return
		}
		ogCtx := r.Context()                                                      // original context
		ctx := context.WithValue(ogCtx, appcontext.UserIDKey, sessionData.UserID) // first er wrap the orginal with a key and its value here it is ID
		ctx = context.WithValue(ctx, appcontext.UserRoleKey, sessionData.Role)    // Second we the first context (ctx) with another key and value (role)
		next.ServeHTTP(w, r.WithContext(ctx))                                     // finally we then pass the ctx with the request which then provides the details to every function.
	})
}

func main() {
	// setting up the connection to the database
	database.ConnectToDB()
	r := mux.NewRouter()

	r.HandleFunc("/signup", handlers.SignUpUser).Methods("POST")
	r.HandleFunc("/login", handlers.LoginUser).Methods("POST")

	// initializing a path prefixx that uses middleware rather then individually typing
	api := r.PathPrefix("/api").Subrouter()
	api.Use(SessionAuthMiddleware)

	// api.Handle("/logout", SessionAuthMiddleware(http.HandlerFunc(handlers.LogoutUser))).Methods("POST")
	// api.Handle("/tickets", SessionAuthMiddleware(http.HandlerFunc(handlers.GetTickets))).Methods("GET")
	api.HandleFunc("/logout", handlers.LogoutUser).Methods("POST")
	api.HandleFunc("/tickets", handlers.GetTickets).Methods("GET")
	api.HandleFunc("/tickets", handlers.CreateTicket).Methods("POST")
	api.HandleFunc("/tickets/{id:[0-9]+}", handlers.DeleteTicket).Methods("DELETE")

	api.HandleFunc("/tickets/assign", handlers.AssignTicket).Methods("POST")
	api.HandleFunc("/tickets/priority", handlers.SetPriority).Methods("POST")
	api.HandleFunc("/tickets/status", handlers.SetStatus).Methods("POST")
	api.HandleFunc("/activity", handlers.ActivityLog).Methods("POST")

	fmt.Println("Server starting on port :8080...")
	log.Fatal(http.ListenAndServe(":8080", r))
}
