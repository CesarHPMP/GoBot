package spotify

import (
	"log"

	"github.com/zmb3/spotify"
	"golang.org/x/oauth2"
)

var client spotify.Client

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
