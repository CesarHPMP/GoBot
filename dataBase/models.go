package database

import (
	"time"
)

type User struct {
	ID           int       `db:"id"`
	SpotifyID    string    `db:"spotify_id"`
	AccessToken  string    `db:"access_token"`
	RefreshToken string    `db:"refresh_token"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

type TopTrack struct {
	ID      int        `db:"id"`
	UserID  int        `db:"user_id"`
	TrackID string     `db:"track_id"`
	Name    string     `db:"name"`
	Artists string     `db:"artists"` // Can store as a comma-separated string
	Album   string     `db:"album"`
	AddedAt time.Month `db:"added_at"`
}
