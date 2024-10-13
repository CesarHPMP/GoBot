package database

import (
	"database/sql"
	"fmt"
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
		user_id TEXT NOT NULL UNIQUE,
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

// SaveTrack saves a user's top tracks to the database
func SaveTrack(userID string, track TopTrack) error {
	// Create tables if they don't exist
	createTables()

	// Find the user's internal ID
	var id int
	err := Database.QueryRow(`SELECT id FROM users WHERE user_id = ?`, userID).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("user not found")
		}
		return err
	}

	// Insert the track data
	query := `INSERT INTO tracks (user_id, track_name, artist, album) VALUES (?, ?, ?, ?)`
	_, err = Database.Exec(query, id, track.Name, track.Artists, track.Album)

	if err != nil {
		return err
	}

	return nil
}
