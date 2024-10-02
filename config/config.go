package config

import (
	"log"
	"os"
)

var (
	BotToken      string
	SpotifyID     string
	SpotifySecret string
)

func LoadConfig() {
	BotToken = os.Getenv("DISCORD_BOT_TOKEN")
	SpotifyID = os.Getenv("SPOTIFY_CLIENT_ID")
	SpotifySecret = os.Getenv("SPOTIFY_CLIENT_SECRET")

	if BotToken == "" || SpotifyID == "" || SpotifySecret == "" {
		log.Fatal("Missing required environment variables.")
	}
}
