package spotify

import {
    "github.com/CesarHPMP/GoBot/database"
}

func (sc *SpotifyClient) StoreTopTracks(userID int) error {
	trackList, err := sc.GetTopTracks()
	if err != nil {
		return err
	}
	
	// Loop through trackList and insert into the database
	for _, track := range trackList {
        db
	}
	return nil
}
