package main

import (
	"log"

	"github.com/CesarHPMP/GoBot/internal/discord"
)

func main() {
	// Create a new instance of the Discord client
	discord.StartBot()

	// You're now connected to Discord and Spotify!
	log.Println("Connected to Discord and Spotify!")
}
