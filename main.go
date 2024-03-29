package main

import (
	"discord-bot/bot"
	"discord-bot/commands"
	db "discord-bot/database"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	// Load the .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Load environment variables
	botToken := os.Getenv("BOT_TOKEN")
	openWeatherToken := os.Getenv("OPENWEATHER_TOKEN")
	translateToken := os.Getenv("TRANSLATE_TOKEN")
	mongoDBToken := os.Getenv("MONGO_DB_TOKEN")
	gptToken := os.Getenv("GPT_TOKEN")

	// Save API keys & start bot
	bot.BotToken = botToken
	bot.TranslateToken = translateToken
	commands.OpenWeatherToken = openWeatherToken
	commands.GPTKey = gptToken
	db.MongoDBURI = mongoDBToken

	bot.Run()
}
