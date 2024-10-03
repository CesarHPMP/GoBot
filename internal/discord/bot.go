package discord

import (
	"log"
	"strings"

	"github.com/CesarHPMP/GoBot/config"
	"github.com/CesarHPMP/GoBot/internal/spotify"
	"github.com/bwmarrin/discordgo"
)

var Finish_run = make(chan bool)

func StartBot() (*discordgo.Session, error) {

	token, _, _ := config.LoadConfig()

	dg, err := discordgo.New("Bot " + token)

	if err != nil {
		return nil, err
	}

	err = dg.Open()
	if err != nil {
		return nil, err
	}
	log.Println("Bot is now running.")

	return dg, nil
}

func MessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore messages from the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Log the received message object
	log.Printf("Received message object: %+v\n", m)

	// Log the raw content of the message
	log.Printf("Message content length: %d\n", len(m.Content))
	log.Printf("Raw message content: %q\n", m.Content)

	// Check if the message content is empty
	if m.Content == "" {
		log.Println("Content is empty")
		return
	}

	// Check for the /connect command
	if strings.HasPrefix(m.Content, "/connect") {
		log.Println("Received /connect command")
		spotify.Starting(s, m.ChannelID)
	}

	// Check for the /turnoff command
	if strings.HasPrefix(m.Content, "/turnoff") {
		log.Println("Received /turnoff command")
		Finish_run <- true
	}
}
