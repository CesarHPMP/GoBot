package main

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/CesarHPMP/GoBot/internal/spotify"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
)

var (
	clientId     string
	clientSecret string
	redirectURI  = "http://localhost:8080/callback"
	initState    string
	config       *oauth2.Config
)

func init() {
	err := godotenv.Load("spotify.env")
	if err != nil {
		log.Fatal(err)
	}
	clientId = os.Getenv("SPOTIFY_CLIENT_ID")
	clientSecret = os.Getenv("SPOTIFY_CLIENT_SECRET")
	config = &oauth2.Config{
		ClientID:     clientId,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURI,
		Scopes:       []string{"user-read-private"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://accounts.spotify.com/authorize",
			TokenURL: "https://accounts.spotify.com/api/token",
		},
	}
}

func main() {
	fmt.Println("Client ID:", clientId)
	http.HandleFunc("/callback", completeAuth)

	initState = generateState()
	if initState == "" {
		log.Fatal("Failed to generate state")
	}
	url := config.AuthCodeURL(initState)
	fmt.Println("Please log in to Spotify by visiting the following page in your browser:", url)

	log.Fatal(http.ListenAndServe(":8080", nil))

	InitSpotify(clientId, clientSecret)
}

func completeAuth(w http.ResponseWriter, r *http.Request) {
	state := r.URL.Query().Get("state")
	if state != initState {
		http.Error(w, "Invalid state", http.StatusBadRequest)
		return
	}

	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Code is missing", http.StatusBadRequest)
		return
	}

	token, err := config.Exchange(context.Background(), code)
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusForbidden)
		log.Fatal(err)
	}

	// Use the token to authenticate the client
	client := config.Client(context.Background(), token)

	// Create a new Spotify client
	spotifyClient := spotify.NewClient(client)

	// Use the Spotify client to make requests to the Spotify API
	user, err := spotifyClient.CurrentUser()
	if err != nil {
		http.Error(w, "Couldn't get user", http.StatusInternalServerError)
		log.Fatal(err)
	}

	fmt.Fprintf(w, "Login Completed! You are logged in as: %s", user.DisplayName)
}

func generateState() string {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatal(err)
	}
	return base64.StdEncoding.EncodeToString(b)
}
