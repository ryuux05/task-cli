package storage

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

const defaultDbFile = "storage/task.db"

// NewSqlite creates a new SQLite database connection to the default database
func NewSqlite() (*sql.DB, error) {
	return connectToSqlite(defaultDbFile)
}

// NewTeamSqlite creates a new SQLite database connection for a specific team
func NewTeamSqlite(teamName string) (*sql.DB, error) {
	// Create the storage directory if it doesn't exist
	err := os.MkdirAll("storage", 0755)
	if err != nil {
		return nil, fmt.Errorf("failed to create storage directory: %v", err)
	}

	dbFile := fmt.Sprintf("storage/team_%s.db", teamName)
	return connectToSqlite(dbFile)
}

// connectToSqlite connects to the specified SQLite database file
func connectToSqlite(dbFile string) (*sql.DB, error) {
	// Ensure the directory exists
	dir := filepath.Dir(dbFile)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create directory %s: %v", dir, err)
	}

	db, err := sql.Open("sqlite", dbFile)
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
