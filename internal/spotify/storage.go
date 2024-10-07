package spotify

import (
	"encoding/csv"
	"errors"
	"os"
	"strconv"
	"strings"
	"time"

	db "github.com/CesarHPMP/GoBot/dataBase"
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

		_, month, year := time.Now().Date()
		date := string(month) + string(year)
		dateint, err := strconv.Atoi(date)

		if err != nil {
			return err
		}

		// Create a Track instance
		track := &db.TopTrack{
			UserID:  0, // You can modify this to reflect the actual user ID
			Name:    trackName,
			Artists: artist,
			Album:   album,
			AddedAt: dateint,
		}

		// Store the track in the database
		if err := db.SaveTrack(track); err != nil {
			return err
		}
	}

	return nil
}
