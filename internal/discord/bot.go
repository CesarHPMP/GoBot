package discord

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

func StartBot(token string) error {
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
