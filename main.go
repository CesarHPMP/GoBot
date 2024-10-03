package main

import (
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"

	"github.com/CesarHPMP/GoBot/internal/spotify"
)

func main() {
	// Load env file
	err := godotenv.Load("../spotify.env")
	if err != nil {
		log.Fatal(err)
	}

	// Create a new instance of the Discord client
	bot_token := os.Getenv("BOT_TOKEN")
	dg, err := discordgo.New("Bot " + bot_token)
	if err != nil {
		log.Fatal(err)
	}

	// Authenticate with the Discord API
	err = dg.Open()
	if err != nil {
		log.Fatal(err)
	}

	// Add a message handler
	dg.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.ID == s.State.User.ID {
			return
		}

		if m.Content == "/connect" || m.Content == "/login" {
			spotify.Starting(s, m.ChannelID)
		}
	})

	// You're now connected to Discord and Spotify!
	log.Println("Connected to Discord and Spotify!")
}
