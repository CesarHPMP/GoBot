package discord

import (
	"log"
	"strings"
	"sync"

	"github.com/CesarHPMP/GoBot/config"
	Myspotify "github.com/CesarHPMP/GoBot/internal/spotify"
	"github.com/bwmarrin/discordgo"
)

var Finish_run = make(chan bool)
var Wg sync.WaitGroup // Add WaitGroup to track async operations

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
	var userID = m.Author.ID

	// Log the received message
	log.Printf("Raw message content: %q\n", m.Content)

	// Check if the message content is empty
	if m.Content == "" {
		log.Println("Content is empty")
		return
	}

	// Local function to check authentication
	checkAuth := func() error {
		if GetSpotifyClient(userID) == nil {
			_, err := s.ChannelMessageSend(m.ChannelID, "Not logged in, please use /connect to log in.")
			return err
		}
		return nil
	}

	// Use switch to handle commands
	switch {
	case strings.HasPrefix(m.Content, "/connect"):
		connected := GetSpotifyClient(m.Author.ID)
		if connected != nil {
			if connected.Connected {
				_, err := s.ChannelMessageSend(m.ChannelID, "You are already logged in.")
				if err != nil {
					return
				}
				return
			}

		}
		userID := m.Author.ID

		Wg.Add(1) // Increment WaitGroup counter for the new goroutine
		go func() {
			defer Wg.Done() // Ensure the counter is decremented when the goroutine finishes

			// Create a new SpotifyClient for the user
			userClient := Myspotify.NewSpotifyClient()
			Myspotify.UserSpotifyClients[userID] = userClient

			// Start the authentication process for this user
			spotifyClient := Myspotify.UserSpotifyClients[userID].Starting(s, m, userID)
			if spotifyClient == nil {
				log.Println("Spotify authentication failed for user:", userID)
				return
			}
		}()

	case strings.HasPrefix(m.Content, "/TopTracks"):
		userID := m.Author.ID

		err := checkAuth() // Pass userID to check authentication
		if err != nil {
			break
		}

		Wg.Add(1) // Increment WaitGroup counter for the new goroutine
		go func() {
			defer Wg.Done()                        // Ensure the counter is decremented when the goroutine finishes
			userClient := GetSpotifyClient(userID) // Get the correct SpotifyClient for this user
			if userClient == nil {
				log.Println("User not authenticated:", userID)
				return
			}

			topTracks, err := Myspotify.UserSpotifyClients[userID].GetTopTracks() // Fetch top tracks for the authenticated user
			if err != nil {
				log.Println("Error fetching top tracks:", err)
				return
			}

			s.ChannelMessageSend(m.ChannelID, topTracks)
		}()

	case strings.HasPrefix(m.Content, "/TopAlbums"):
		if err := checkAuth(); err != nil {
			log.Println("Error checking authentication:", err)
			break
		}

		Wg.Add(1) // Increment WaitGroup counter for the new goroutine
		go func() {
			defer Wg.Done() // Ensure the counter is decremented when the goroutine finishes
			topAlbums, err := Myspotify.UserSpotifyClients[userID].GetTopAlbums()
			if err != nil {
				log.Println(err)
				return
			}
			s.ChannelMessageSend(m.ChannelID, topAlbums)
		}()
	}
}

func GetSpotifyClient(userID string) *Myspotify.SpotifyClient {
	if user_client, exists := Myspotify.UserSpotifyClients[userID]; exists {
		return user_client
	} else {
		log.Print("User client does not exist")
		return nil
	}
}
