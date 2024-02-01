package bot

import (
	"discord-bot/commands"
	"strings"

	"github.com/bwmarrin/discordgo"
)
func newMessage(discord *discordgo.Session, message *discordgo.MessageCreate) {
	// Ignore bot message
	if message.Author.ID == discord.State.User.ID {
		return
	}
	// Respond to messages
	switch {
	case strings.Contains(message.Content, "!weather"):
		currentWeather := commands.GetCurrentWeather(message.Content)
		discord.ChannelMessageSendComplex(message.ChannelID, currentWeather)
	case strings.Contains(message.Content, "!translate"):
		translation := commands.TranslateText(message.Content, TranslateToken)
		discord.ChannelMessageSendComplex(message.ChannelID, translation)
	case strings.HasPrefix(message.Content, "!remind"):
		commands.HandleRemindCommand(discord, message)	
	case strings.HasPrefix(message.Content, "!gpt"):
		commands.HandleGPTCommand(discord, message)
	case strings.HasPrefix(message.Content, "!rps"):
		commands.HandleRPSCommand(discord, message)
	case strings.HasPrefix(message.Content, "!poll"):
		commands.HandlePollCommand(discord, message)
	case strings.HasPrefix(message.Content, "!vote"):
		commands.HandleVoteCommand(discord, message)
	case strings.HasPrefix(message.Content, "!result"):
		commands.HandleResultCommand(discord, message)
	case strings.HasPrefix(message.Content, "!help"):
		commands.HandleHelpCommand(discord, message)
	}
	
}