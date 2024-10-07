package spotify

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"github.com/zmb3/spotify"
	"golang.org/x/oauth2"

	Setup "github.com/CesarHPMP/GoBot/config"
	"github.com/CesarHPMP/GoBot/utils"
)

// SpotifyClient struct to store individual user session data
type SpotifyClient struct {
	Client      *spotify.Client
	Connected   bool
	AuthDone    chan bool // Channel to signal when auth is done
	Config      *oauth2.Config
	RedirectURI string
	State       string
}

var (
	clientId     string
	clientSecret string
	defaultURI   = "https://select-sheep-currently.ngrok-free.app/callback" // Default URI
)

// Initialize the default Spotify settings (env load)
func init() {
	err := godotenv.Load("../spotify.env")
	if err != nil {
		log.Fatal(err)
	}

	_, clientId, clientSecret = Setup.LoadConfig()
}

// NewSpotifyClient creates a new SpotifyClient instance for a user
func NewSpotifyClient() *SpotifyClient {
	config := &oauth2.Config{
		ClientID:     clientId,
		ClientSecret: clientSecret,
		RedirectURL:  defaultURI,
		Scopes:       []string{"user-read-private", "user-top-read"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://accounts.spotify.com/authorize",
			TokenURL: "https://accounts.spotify.com/api/token",
		},
	}

	return &SpotifyClient{
		Config:      config,
		AuthDone:    make(chan bool), // Each user has their own channel for signaling auth completion
		RedirectURI: defaultURI,
		State:       utils.GenerateState(),
	}
}

// Starting initiates the Spotify login process for a specific user
func (sc *SpotifyClient) Starting(dg *discordgo.Session, m *discordgo.MessageCreate, userID string) *spotify.Client {
	if dg.Client == nil {
		log.Fatal("Discord session client is nil")
	}

	url := sc.Config.AuthCodeURL(sc.State)
	_, err := dg.ChannelMessageSend(m.ChannelID, "Please log in to Spotify by visiting the following page in your browser: "+url)
	if err != nil {
		log.Fatal(err)
	}

	var port = ":8080"

	go func() {
		http.HandleFunc("/callback", sc.CompleteAuth)
		if err := http.ListenAndServe(port, nil); err != nil && err != http.ErrServerClosed {
			log.Println("HTTP server error:", err)
			srv := http.Server{Addr: port}
			srv.Close()
			return
		}
	}()

	srv := &http.Server{Addr: port}

	// Create a timeout channel to avoid waiting indefinitely
	timeout := time.After(30 * time.Second)

	// Wait for either the auth process to finish or the timeout
	select {
	case <-sc.AuthDone: // Authentication successful
		fmt.Println("User logged in successfully!")
		_, _ = dg.ChannelMessageSend(m.ChannelID, "Authentication successful! You are now logged in.")
		sc.Connected = true
		srv.Close() // Close server after auth
		return sc.Client

	case <-timeout: // Timeout occurred
		fmt.Println("Authentication timed out after 30 seconds.")
		_, _ = dg.ChannelMessageSend(m.ChannelID, "Authentication timed out. Please try again.")
		srv.Close() // Close server after timeout
		return nil
	}
}

func (sc *SpotifyClient) CompleteAuth(w http.ResponseWriter, r *http.Request) {
	state := r.URL.Query().Get("state")
	if state != sc.State {
		http.Error(w, "Invalid state", http.StatusBadRequest)
		return
	}

	// Gets code for auth
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Code is missing", http.StatusBadRequest)
		return
	}

	token, err := sc.Config.Exchange(context.Background(), code)
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusForbidden)
		log.Fatal(err)
	}

	sc.InitSpotify(token)
	if sc.Client == nil {
		http.Error(w, "Couldn't get client", http.StatusForbidden)
		return
	}

	// Notify the main thread that authentication is complete
	sc.AuthDone <- true
}

func (sc *SpotifyClient) InitSpotify(token *oauth2.Token) {
	c := spotify.Authenticator{}.NewClient(token)
	sc.Client = &c
}

// GetTopTracks returns the top tracks for the authenticated user
func (sc *SpotifyClient) GetTopTracks() (string, error) {
	if sc.Client == nil {
		return "", errors.New("client is not initialized")
	}

	topTracks, err := sc.Client.CurrentUsersTopTracks()
	if err != nil {
		return "", err
	}

	var trackList string
	for i, track := range topTracks.Tracks {
		artistNames := make([]string, len(track.Artists))
		for i, artist := range track.Artists {
			artistNames[i] = artist.Name
		}
		trackList += fmt.Sprintf("%d. %s - %s (Album: %s)\n", i+1, track.Name, strings.Join(artistNames, ", "), track.Album.Name)
	}
	return trackList, nil
}

// GetTopAlbums returns the top albums for the authenticated user
func (sc *SpotifyClient) GetTopAlbums() (string, error) {
	if sc.Client == nil {
		return "", errors.New("client is not initialized")
	}

	type albumCount struct {
		Album   string
		Count   int
		Artists string
	}

	limit := 50
	timeRange := "long" // Use Spotify's time range keyword
	hashTable := utils.NewHashTable()

	var options = &spotify.Options{
		Limit:     &limit,
		Timerange: &timeRange,
	}

	topTracks, err := sc.Client.CurrentUsersTopTracksOpt(options)
	if err != nil {
		return "", err
	}

	var albums []albumCount

	for _, track := range topTracks.Tracks {
		count := hashTable.Get(track.Album.Name)
		alreadyInList := false

		// Check if album already exists in the slice
		for _, a := range albums {
			if a.Album == track.Album.Name {
				alreadyInList = true
				break
			}
		}

		// Only add the album if it hasn't been added yet
		if !alreadyInList {
			// Get all artists' names and join them in a string
			artistNames := make([]string, len(track.Album.Artists))
			for i, artist := range track.Album.Artists {
				artistNames[i] = artist.Name
			}

			albums = append(albums, albumCount{
				Album:   track.Album.Name,
				Count:   count,
				Artists: strings.Join(artistNames, ", "),
			})
		}
	}

	// Sort albums by count in descending order
	sort.Slice(albums, func(i, j int) bool {
		return albums[i].Count > albums[j].Count
	})

	// Prepare a string to hold the formatted album list
	albumList := ""
	for i, album := range albums {
		albumList += fmt.Sprintf("%d. %s by %s\n", i+1, album.Album, album.Artists)
		if i == 9 { // Limit to the top 10 albums
			break
		}
	}

	// Return the formatted string containing the top albums without duplicates
	return albumList, nil
}
