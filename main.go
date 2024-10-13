package main

import (
	"log"

	database "github.com/CesarHPMP/GoBot/dataBase"
	"github.com/CesarHPMP/GoBot/internal/discord"
)

func main() {
	// Initialize database
	dataSource := database.GetDataSourceName()
	database.InitDB(dataSource)
	// Start the bot
	dg, err := discord.StartBot()
	if err != nil {
		log.Fatalf("Error starting bot: %v", err)
	}

	dg.AddHandler(discord.MessageCreate)

	<-discord.Finish_run
	log.Println("Bot shutting down due to /turnoff command.")

	// Wait for all goroutines (async tasks) to finish
	discord.Wg.Wait()

	// Close the bot
	if err = dg.Close(); err != nil {
		log.Fatalf("Error closing bot: %v", err)
	}
}
