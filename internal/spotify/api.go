package spotify

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/zmb3/spotify"
	"golang.org/x/oauth2"

	"github.com/CesarHPMP/GoBot/utils"
	"github.com/bwmarrin/discordgo"
)

var client spotify.Client

var (
	clientId     string
	clientSecret string
	redirectURI  = "http://localhost:8080/callback"
	initState    string
	config       *oauth2.Config
)

func init() {

	err := godotenv.Load("../spotify.env")
	if err != nil {
		log.Fatal(err)
	}

	clientId := os.Getenv("SPOTIFY_CLIENT_ID")
	clientSecret := os.Getenv("SPOTIFY_CLIENT_SECRET")
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

}

func completeAuth(w http.ResponseWriter, r *http.Request) {
	//gets state and compares it to the one generated
	state := r.URL.Query().Get("state")
	if state != initState {
		http.Error(w, "Invalid state", http.StatusBadRequest)
		return
	}

	// gets code for auth
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
	topTracks, err := GetTopTracks()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(topTracks)
}

func InitSpotify(config *oauth2.Config, token *oauth2.Token) {
	client = spotify.Authenticator{}.NewClient(token)
}

func GetTopTracks() (*spotify.FullTrackPage, error) {
	topTracks, err := client.CurrentUsersTopTracks()
	if err != nil {
		log.Fatal("Error getting top tracks:", err)
		return nil, err
	}

	return topTracks, nil
}
