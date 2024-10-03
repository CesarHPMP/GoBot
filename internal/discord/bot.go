package discord

import (
	"log"

	"github.com/CesarHPMP/GoBot/config"
	"github.com/CesarHPMP/GoBot/internal/spotify"
	"github.com/bwmarrin/discordgo"
)

func StartBot() error {

	token, _, _ := config.LoadConfig()

	dg, err := discordgo.New("Bot " + token)

	if err != nil {
		return err
	}

	dg.AddHandler(messageCreate)

	err = dg.Open()
	if err != nil {
		return err
	}

	log.Println("Bot is now running.")
	return nil
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Content == "!ping" {
		s.ChannelMessageSend(m.ChannelID, "Pong!")
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
