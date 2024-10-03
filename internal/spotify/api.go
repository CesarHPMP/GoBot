package spotify

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/zmb3/spotify"
	"golang.org/x/oauth2"

	Setup "github.com/CesarHPMP/GoBot/config"
	"github.com/CesarHPMP/GoBot/utils"
	"github.com/bwmarrin/discordgo"
)

var Client spotify.Client

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

	// Wait for the auth process to finish
	fmt.Println("Waiting for user to log in...")
	<-authDone // Block until the auth is complete
	fmt.Println("User logged in successfully! Continuing with the flow...")

	top_tracks, err := GetTopTracks()
	dg.ChannelMessageSend(channelID, top_tracks)

	if err != nil {
		log.Fatal(err)
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

	// Notify the main thread that authentication is complete
	authDone <- true // Signal that authentication is finished
}

func InitSpotify(config *oauth2.Config, token *oauth2.Token) {
	Client = spotify.Authenticator{}.NewClient(token)
}

func GetTopTracks() (string, error) {
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
