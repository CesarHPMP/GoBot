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

	if s.Client == nil {
		return
	}

	// Ignore messages from the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Log the received message object
	log.Printf("Received message object: %+v\n", m)
	log.Printf("Message content length: %d\n", len(m.Content))
	log.Printf("Raw message content: %q\n", m.Content)

	// Check if the message content is empty
	if m.Content == "" {
		log.Println("Content is empty")
		return
	}

	// Local function to check authentication
	checkAuth := func() error {
		if !spotify.Connected {
			_, err := s.ChannelMessageSend(m.ChannelID, "Not logged in, please use /connect to log in.")
			return err
		}
		return nil
	}

	// Use switch to handle commands
	switch {
	case strings.HasPrefix(m.Content, "/connect"):
		// Connect to Spotify
		spotify.Starting(s, m.ChannelID)

	case strings.HasPrefix(m.Content, "/TopTracks"):
		err := checkAuth()
		if err != nil {
			break
		}
		topTracks, err := spotify.GetTopTracks()
		if err != nil {
			log.Println(err)
			return
		}
		s.ChannelMessageSend(m.ChannelID, topTracks)

	case strings.HasPrefix(m.Content, "/TopAlbums"):
		err := checkAuth()
		if err != nil {
			break
		}
		topAlbums, err := spotify.GetTopAlbums()
		if err != nil {
			log.Println(err)
			return
		}
		s.ChannelMessageSend(m.ChannelID, topAlbums)

	case strings.HasPrefix(m.Content, "/turnoff"):
		log.Println("Received /turnoff command")
		Finish_run <- true
	}
}
