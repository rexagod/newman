package database

import (
	"database/sql"
	"log"
	"os"
	"strings"

	_ "github.com/lib/pq"
)

var db *sql.DB

// Initialize connects to the database and pings it.
func Initialize() error {
	connStr := os.Getenv("DB_CONN_STRING")

	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		return err
	}

	err = db.Ping()
	if err != nil {
		return err
	}

	log.Println("Connected to yo database!")
	return nil
}

// GetDB returns the database connection.
func GetDB() *sql.DB {
	return db
}

// ExecuteQueriesFromFile reads and executes queries from a SQL file.
func CreateTables() error {
	content, err := os.ReadFile("database/tables.sql")

	if err != nil {
		return err
	}
	queries := strings.Split(string(content), ";")

	for _, query := range queries {
		_, err := db.Exec(query)
		if err != nil {
			return err
		}
	}

	log.Println("Executed queries from tables.sql.")
	return nil
}
