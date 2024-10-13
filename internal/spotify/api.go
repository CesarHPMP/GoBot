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
	"github.com/gorilla/mux"
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
	clientId           string
	clientSecret       string
	defaultURI         = "https://select-sheep-currently.ngrok-free.app/callback" // Default URI
	UserSpotifyClients = make(map[string]*SpotifyClient)
)

// Initialize the default Spotify settings (env load)
func init() {
	if err := godotenv.Load("../spotify.env"); err != nil {
		log.Fatal("Error loading .env file:", err)
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

	// Generate the authorization URL with userID in the state parameter
	url := sc.Config.AuthCodeURL(sc.State + ":" + userID)

	// Send the authorization URL to the user
	if _, err := dg.ChannelMessageSend(m.ChannelID, "Please log in to Spotify by visiting the following page in your browser: "+url); err != nil {
		log.Println("Error sending message:", err)
	}

	port := ":8080"
	router := mux.NewRouter()

	// Register the callback handler
	router.HandleFunc("/callback/{userID}", sc.CompleteAuth).Methods("GET")

	srv := &http.Server{
		Addr:    port,
		Handler: router,
	}

	// Create a timeout channel to avoid waiting indefinitely
	timeout := time.After(30 * time.Second)

	go func() {
		// Start the HTTP server
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Println("HTTP server error:", err)
		}
	}()

	// Wait for either the auth process to finish or the timeout
	select {
	case <-sc.AuthDone: // Authentication successful
		log.Println("User logged in successfully!")
		if _, err := dg.ChannelMessageSend(m.ChannelID, "Authentication successful! You are now logged in."); err != nil {
			log.Println("Error sending success message:", err)
		}
		sc.Connected = true

		// Close the server after auth
		if err := srv.Shutdown(context.Background()); err != nil {
			log.Println("Server Shutdown Failed:", err)
		}
		return sc.Client

	case <-timeout: // Timeout occurred
		log.Println("Authentication timed out after 30 seconds.")
		if _, err := dg.ChannelMessageSend(m.ChannelID, "Authentication timed out. Please try again."); err != nil {
			log.Println("Error sending timeout message:", err)
		}
		if err := srv.Shutdown(context.Background()); err != nil {
			log.Println("Server Shutdown Failed:", err)
		}
		return nil
	}
}

// CompleteAuth handles the Spotify callback and completes the authentication process
func (sc *SpotifyClient) CompleteAuth(w http.ResponseWriter, r *http.Request) {
	// Extract the state parameter
	state := r.URL.Query().Get("state")
	parts := strings.Split(state, ":")
	if len(parts) != 2 {
		http.Error(w, "Invalid state", http.StatusBadRequest)
		log.Println("Invalid state parameter:", state)
		return
	}

	userID := parts[1] // Extract userID from the state parameter

	// Validate the state
	if state != sc.State {
		http.Error(w, "Invalid state", http.StatusBadRequest)
		log.Println("State mismatch:", state)
		return
	}

	// Get the authorization code from the callback
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Code is missing", http.StatusBadRequest)
		log.Println("Authorization code is missing")
		return
	}

	// Exchange the code for a token
	token, err := sc.Config.Exchange(context.Background(), code)
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusForbidden)
		log.Println("Token exchange error:", err)
		return
	}

	// Initialize the Spotify client
	sc.InitSpotify(token)
	if sc.Client == nil {
		http.Error(w, "Couldn't get client", http.StatusForbidden)
		log.Println("Spotify client initialization failed")
		return
	}

	// Notify the main thread that authentication is complete
	sc.AuthDone <- true

	// Optionally send a success response back to the user
	_, _ = fmt.Fprintf(w, "Authentication successful for user %s! You can close this window.", userID)
}

// InitSpotify initializes the Spotify client with the provided token
func (sc *SpotifyClient) InitSpotify(token *oauth2.Token) {
	c := spotify.Authenticator{}.NewClient(token)
	sc.Client = &c
}

// GetTopTracks returns the top tracks for the authenticated user
func (sc *SpotifyClient) GetTopTracks() (string, error) {
	if sc.Client == nil {
		err := errors.New("client is not initialized")
		log.Println("GetTopTracks Error:", err) // Log the error
		return "", err
	}

	topTracks, err := sc.Client.CurrentUsersTopTracks()
	if err != nil {
		log.Println("Failed to get top tracks:", err) // Log the error
		return "", err
	}

	var trackList strings.Builder // Use a strings.Builder for better performance
	for i, track := range topTracks.Tracks {
		artistNames := make([]string, len(track.Artists))
		for i, artist := range track.Artists {
			artistNames[i] = artist.Name
		}
		trackList.WriteString(fmt.Sprintf("%d. %s - %s (Album: %s)\n", i+1, track.Name, strings.Join(artistNames, "; "), track.Album.Name))
	}

	userID, exists := GetUserIDFromClient(sc, UserSpotifyClients)
	if !exists {
		return "", errors.New("user ID not found")
	}

	InputTopTracks(trackList.String(), userID)

	return trackList.String(), nil
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

	options := &spotify.Options{
		Limit:     &limit,
		Timerange: &timeRange,
	}

	topTracks, err := sc.Client.CurrentUsersTopTracksOpt(options)
	if err != nil {
		log.Println("Failed to get top tracks for albums:", err)
		return "", err
	}

	albums := []albumCount{}

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

	// Prepare a string to display the top albums
	var albumList strings.Builder
	for i, album := range albums {
		albumList.WriteString(fmt.Sprintf("%d. %s by %s - Count: %d\n", i+1, album.Album, album.Artists, album.Count))
	}

	return albumList.String(), nil
}

// GetUserIDFromClient retrieves the user ID associated with the SpotifyClient
func GetUserIDFromClient(sc *SpotifyClient, userMap map[string]*SpotifyClient) (string, bool) {
	for userID, client := range userMap {
		if client == sc {
			return userID, true
		}
	}
	return "", false
}
