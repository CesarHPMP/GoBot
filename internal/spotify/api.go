package spotify

import (
	"log"

	"github.com/zmb3/spotify"
	"golang.org/x/oauth2"
)

var client spotify.Client

<<<<<<< HEAD
func InitSpotify(config *oauth2.Config, token *oauth2.Token) {
	client = spotify.Authenticator{}.NewClient(token)
=======
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
>>>>>>> origin/revert-1-initial_dev
}

func GetTopTracks() (*spotify.FullTrackPage, error) {
	topTracks, err := client.CurrentUsersTopTracks()
	if err != nil {
		log.Fatal("Error getting top tracks:", err)
		return nil, err
	}

	return topTracks, nil
}
