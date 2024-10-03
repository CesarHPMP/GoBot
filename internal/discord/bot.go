package discord

import (
	"fmt"
	"log"
	"strings"

	"github.com/CesarHPMP/GoBot/config"
	"github.com/CesarHPMP/GoBot/internal/spotify"
	"github.com/bwmarrin/discordgo"
)

func StartBot() (*discordgo.Session, error) {

	token, _, _ := config.LoadConfig()

	dg, err := discordgo.New("Bot " + token)

	if err != nil {
		return nil, err
	}

	dg.AddHandler(ReadMessageAndConnect)

	err = dg.Open()
	if err != nil {
		return nil, err
	}

	ChannelID := "1291147572265746524"

	// Send a message to the specific channel
	_, err = dg.ChannelMessageSend(ChannelID, "Hello, Discord channel!")
	if err != nil {
		fmt.Println("Error sending message,", err)
	}

	log.Println("Bot is now running.")
	return dg, nil
}

func MessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if strings.HasPrefix(m.Content, "/connect") {
		spotify.Starting(s, m.ChannelID)
	}
}

func ReadMessageAndConnect(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Content == "/connect" {
		spotify.Starting(s, m.ChannelID)
	}
}
