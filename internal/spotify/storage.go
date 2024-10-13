package spotify

import (
	"strconv"
	"strings"
	"time"

	db "github.com/CesarHPMP/GoBot/dataBase"
)

// InputTopTracks reads top tracks from a specified CSV file and stores them in the database
func InputTopTracks(trackpack string, userID string) error {

	tracks := strings.Split(trackpack, ". -")

	for _, record := range tracks {

		trackName := string(record[0])
		artist := string(record[1])
		album := string(record[2])

		date := time.Now().Month()
		sumName, err := strconv.Atoi(trackName)

		if err != nil {
			return err
		}

		sumArtist, err := strconv.Atoi(artist)

		if err != nil {
			return err
		}

		sumAlbum, err := strconv.Atoi(album)

		if err != nil {
			return err
		}

		trackSum := sumAlbum + sumArtist + sumName
		trackID := strconv.Itoa(trackSum)

		intUserID, err := strconv.Atoi(userID)

		if err != nil {
			return err
		}

		// Create a Track instance
		track := &db.TopTrack{
			ID:      1 + 1,
			UserID:  intUserID, // You can modify this to reflect the actual user ID
			TrackID: trackID,
			Name:    trackName,
			Artists: artist,
			Album:   album,
			AddedAt: date,
		}

		// Store the track in the database
		if err := db.SaveTrack(userID, *track); err != nil {
			return err
		}
	}

	return nil
}
