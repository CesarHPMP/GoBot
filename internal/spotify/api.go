package spotify

import (
	"context"
	"log"

	"github.com/zmb3/spotify"
	"golang.org/x/oauth2/clientcredentials"
)

var client *spotify.Client

func InitSpotify(clientID, clientSecret string) {
	config := &clientcredentials.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		TokenURL:     spotify.TokenURL,
	}
	token, err := config.Token(context.Background())
	if err != nil {
		log.Fatal("Error getting Spotify token:", err)
	}

	spotifyClient := spotify.Authenticator{}.NewClient(token)
}

func GetTopTracks() {
	user, err := client.CurrentUser()
	if err != nil {
		log.Fatal("Error getting current user:", err)
	}

	topTracks, err := client.CurrentUsersTopTracks()
	if err != nil {
		log.Fatal("Error getting top tracks:", err)
	}

	log.Printf("Top tracks for user %s: %+v", user.DisplayName, topTracks)
}
