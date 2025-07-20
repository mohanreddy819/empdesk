package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func ConnectToDB() {
	var err error
	DB, err = sql.Open("mysql", "root:2003@tcp(127.0.0.1:3306)/empdesk")
	if err != nil {
		log.Fatal("Error preparing database connection:", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatal("Error connecting to database:", err)
	}

	fmt.Println("database connected..")

	CreateSchema()
}

func CreateSchema() {

	UserTable := `CREATE TABLE IF NOT EXISTS users (
		id INT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(48) NOT NULL,
		email VARCHAR(255) NOT NULL UNIQUE,
		password TEXT NOT NULL,
		role VARCHAR(15) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`
	if _, err := DB.Exec(UserTable); err != nil {
		log.Fatal("Error creating users table:", err)
	}

	TicketCategories := `CREATE TABLE IF NOT EXISTS ticket_categories (
		id INT AUTO_INCREMENT PRIMARY KEY,
		name TEXT NOT NULL
	);`
	if _, err := DB.Exec(TicketCategories); err != nil {
		log.Fatal("Error creating ticket_categories table:", err)
	}

	Ticket := `CREATE TABLE IF NOT EXISTS tickets (
		id INT AUTO_INCREMENT PRIMARY KEY,
		ticket_token INTEGER NOT NULL UNIQUE,
		user_id INT NOT NULL,
		category_id INT NOT NULL,
		title TEXT NOT NULL,
		description TEXT NOT NULL,
		status VARCHAR(20) DEFAULT 'Open',
		priority VARCHAR(20),
		assigned_to VARCHAR(255),
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
		FOREIGN KEY (category_id) REFERENCES ticket_categories(id) ON DELETE CASCADE
	);`
	if _, err := DB.Exec(Ticket); err != nil {
		log.Fatal("Error creating tickets table:", err)
	}

	TicketComment := `CREATE TABLE IF NOT EXISTS ticket_comments (
		id INT AUTO_INCREMENT PRIMARY KEY,
		ticket_id INT NOT NULL,
		user_id INT NOT NULL,
		comments TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (ticket_id) REFERENCES tickets(id) ON DELETE CASCADE,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	);`
	if _, err := DB.Exec(TicketComment); err != nil {
		log.Fatal("Error creating ticket_comments table:", err)
	}

	Activity := `CREATE TABLE IF NOT EXISTS activity_logs (
		id INT AUTO_INCREMENT PRIMARY KEY,
		user_id INT,
		action TEXT,
		timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL
	);`
	if _, err := DB.Exec(Activity); err != nil {
		log.Fatal("Error creating activity_logs table:", err)
	}
}
