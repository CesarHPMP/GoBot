package config

import (
	"log"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

var (
	BotToken      string
	SpotifyID     string
	SpotifySecret string
)

var Router = mux.NewRouter()

func LoadConfig() (string, string, string) {
	err := godotenv.Load("../spotify.env")

	if err != nil {
		log.Fatal(err)
	}

	BotToken = os.Getenv("BOT_TOKEN")
	SpotifyID = os.Getenv("SPOTIFY_CLIENT_ID")
	SpotifySecret = os.Getenv("SPOTIFY_CLIENT_SECRET")

	if BotToken == "" || SpotifyID == "" || SpotifySecret == "" {
		log.Fatal("Missing required environment variables.")
	}
	return BotToken, SpotifyID, SpotifySecret
}
