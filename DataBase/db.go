package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

var Database *sql.DB

// InitDB initializes the database connection
func InitDB(dataSourceName string) {
	var err error
	Database, err = sql.Open("sqlite3", dataSourceName)
	if err != nil {
		log.Fatal(err)
	}

	// Create tables if they don't exist
	createTables()
}

// createTables creates the necessary tables for storing user data and tracks
func createTables() {
	// Users table
	createUserTable := `CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id TEXT NOT NULL UNIQUE,
		access_token TEXT,
		refresh_token TEXT,
		expires_at DATETIME
	);`

	if _, err := Database.Exec(createUserTable); err != nil {
		log.Fatal(err)
	}

	// Tracks table
	createTrackTable := `CREATE TABLE IF NOT EXISTS tracks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER,
		track_name TEXT,
		artist TEXT,
		album TEXT,
		timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY(user_id) REFERENCES users(id)
	);`

	if _, err := Database.Exec(createTrackTable); err != nil {
		log.Fatal(err)
	}
}
