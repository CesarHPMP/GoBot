package spotify

import (
	"encoding/csv"
	"errors"
	"os"
	"strings"

	"github.com/CesarHPMP/GoBot/database"
)

// InputTopTracks reads top tracks from a specified CSV file and stores them in the database
func InputTopTracks(dataSourceName string) error {
	file, err := os.Open(dataSourceName)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	// Optionally, you can configure the reader (like setting the delimiter)
	reader.Comma = ',' // Assuming the CSV is comma-delimited

	// Read all records from the CSV
	records, err := reader.ReadAll()
	if err != nil {
		return err
	}

	// Loop through the records and save them to the database
	for _, record := range records {
		// Assuming the CSV format: Track Name, Artist, Album
		if len(record) < 3 {
			return errors.New("record has insufficient fields: " + strings.Join(record, ", "))
		}

		trackName := strings.TrimSpace(record[0])
		artist := strings.TrimSpace(record[1])
		album := strings.TrimSpace(record[2])

		// Create a Track instance
		track := &database.Track{
			UserID:    0, // You can modify this to reflect the actual user ID
			TrackName: trackName,
			Artist:    artist,
			Album:     album,
		}

		// Store the track in the database
		if err := database.SaveTrack(track); err != nil {
			return err
		}
	}

	return nil
}
