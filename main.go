package main

import (
	"log"

	"github.com/CesarHPMP/GoBot/internal/discord"
	"github.com/CesarHPMP/GoBot/internal/spotify"
)

func main() {
	// Create a new instance of the Discord client and store it in dg
	dg, err := discord.StartBot()

	channelID := "1291147572265746524"

	if err != nil {
		log.Fatal(err)
	}

	spotify.Starting(dg, "1291147572265746524")

	top_tracks, err := spotify.GetTopTracks()

	if err != nil {
		log.Fatal(err)
	}

	dg.ChannelMessageSend(channelID, top_tracks)

	// You're now connected to Discord and Spotify!
	log.Println("Connected to Discord and Spotify!")
}
