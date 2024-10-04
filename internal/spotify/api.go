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

	"github.com/joho/godotenv"
	"github.com/zmb3/spotify"
	"golang.org/x/oauth2"

	Setup "github.com/CesarHPMP/GoBot/config"
	"github.com/CesarHPMP/GoBot/utils"
	"github.com/bwmarrin/discordgo"
)

var Client *spotify.Client
var Connected bool

var (
	clientId     string
	clientSecret string
	redirectURI  = "https://select-sheep-currently.ngrok-free.app/callback"
	initState    string
	config       *oauth2.Config
	authDone     = make(chan bool) // Channel to signal when auth is done
)

func init() {

	err := godotenv.Load("../spotify.env")
	if err != nil {
		log.Fatal(err)
	}

	_, clientId, clientSecret = Setup.LoadConfig()

	config = &oauth2.Config{
		ClientID:     clientId,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURI,
		Scopes:       []string{"user-read-private", "user-top-read"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://accounts.spotify.com/authorize",
			TokenURL: "https://accounts.spotify.com/api/token",
		},
	}
}

func Starting(dg *discordgo.Session, channelID string) {
	if dg.Client == nil {
		log.Fatal("Client is nil")
	}

	fmt.Println("Client ID:", clientId)
	http.HandleFunc("/callback", completeAuth)

	initState = utils.GenerateState()

	if initState == "" {
		log.Fatal("Failed to generate state")
	}

	url := config.AuthCodeURL(initState)
	_, err := dg.ChannelMessageSend(channelID, "Please log in to Spotify by visiting the following page in your browser (30 seconds limit before process is killed): "+url)

	if err != nil {
		log.Fatal(err)
	}

	go func() {
		log.Fatal(http.ListenAndServe(":8080", nil))
	}()

	// Create a timeout channel that will signal after 30 seconds
	timeout := time.After(30 * time.Second)

	// Wait for either the auth process to finish or the timeout to occur
	select {
	case <-authDone: // Authentication was successful
		fmt.Println("User logged in successfully! Continuing with the flow...")
		_, _ = dg.ChannelMessageSend(channelID, "Authentication successful! You are now logged in.")
		Connected = true
	case <-timeout: // Timeout occurred
		fmt.Println("Authentication timed out after 30 seconds.")
		_, _ = dg.ChannelMessageSend(channelID, "Authentication timed out. Please try again.")
		// Cancel the auth process by closing the HTTP server
		go func() {
			http.DefaultServeMux = new(http.ServeMux) // Reset the default serve mux to stop listening
		}()
	}
}

func completeAuth(w http.ResponseWriter, r *http.Request) {

	// Gets state and compares it to the one generated
	state := r.URL.Query().Get("state")
	if state != initState {
		http.Error(w, "Invalid state", http.StatusBadRequest)
		return
	}

	// Gets code for auth
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Code is missing", http.StatusBadRequest)
		return
	}

	// Gets token for client
	token, err := config.Exchange(context.Background(), code)
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusForbidden)
		log.Fatal(err)
	}

	InitSpotify(config, token)

	if Client == nil {
		http.Error(w, "Couldn't get client", http.StatusForbidden)
		return
	}

	// Notify the main thread that authentication is complete
	authDone <- true // Signal that authentication is finished
}

func InitSpotify(config *oauth2.Config, token *oauth2.Token) {
	c := spotify.Authenticator{}.NewClient(token)
	Client = &c
}

func GetTopTracks() (string, error) {
	if Client == nil {
		return "", errors.New("client is not initialized")
	}

	// Fetch the top tracks
	topTracks, err := Client.CurrentUsersTopTracks()
	if err != nil {
		log.Fatal("Error getting top tracks:", err)
		return "", err
	}

	// Prepare a string to hold the formatted track info
	var trackList string

	// Iterate over the top tracks and extract necessary info
	for i, track := range topTracks.Tracks {
		// Each track has a list of artists, album, and name, which you can format
		artistNames := ""
		for _, artist := range track.Artists {
			artistNames += artist.Name + ", "
		}
		// Remove the trailing comma and space
		artistNames = artistNames[:len(artistNames)-2]

		// Append each track's info to the trackList string
		trackList += fmt.Sprintf("%d. %s - %s (Album: %s)\n", i+1, track.Name, artistNames, track.Album.Name)
	}

	// Return the formatted string containing the top tracks
	return trackList, nil
}

func GetTopAlbums() (string, error) {
	if Client == nil {
		return "", errors.New("client is not initialized")
	}

	limit := 50
	timeRange := "medium" // Use Spotify's time range keyword
	hashTable := utils.NewHashTable()

	var options = &spotify.Options{
		Limit:     &limit,
		Timerange: &timeRange,
	}

	// Fetch the top 50 tracks from the user’s account
	topTracks, err := Client.CurrentUsersTopTracksOpt(options)
	if err != nil {
		return "", err
	}

	// Count how many times each album appears, without repeating albums
	for _, track := range topTracks.Tracks {
		hashTable.Add(track.Album.Name) // Track album appearances
	}

	// Extract albums and their counts into a slice of key-value pairs
	type albumCount struct {
		Album   string
		Count   int
		Artists string
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
