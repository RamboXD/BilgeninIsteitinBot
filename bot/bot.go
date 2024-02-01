package bot

import (
	"discord-bot/commands"
	db "discord-bot/database"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
)

// Store Bot API Tokens:
var (
	OpenWeatherToken string
	BotToken         string
	TranslateToken	 string
	MongoDBURI       string
	GPTKey 			 string
)

func Run() {
	// Create new Discord Session
	db.InitMongoDB()

	discord, err := discordgo.New("Bot " + BotToken)
	// fmt.Println(discord)
	if err != nil {
		log.Fatal(err)
	}

	// Add event handler for general messages
	discord.AddHandler(newMessage)

	go commands.RemindScheduler(discord)


	// Open session
	if err := discord.Open(); err != nil {
		log.Fatalf("Error opening Discord session: %v", err)
	}

	defer discord.Close()

	// Run until code is terminated
	fmt.Println("Bot running...")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

}
