// database.go
package util

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func InitDB() (*sql.DB, error) {
	// Connect to PostgreSQL database
	db, err := sql.Open("postgres", "port=5433 user=trustedauthority password=trustedauthority dbname=trustedauthority sslmode=disable")
	if err != nil {
		return nil, err
	}

	// Check if connection is successful
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	log.Println("Connected to PostgreSQL database")
	return db, nil
}
