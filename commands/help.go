package commands

import "github.com/bwmarrin/discordgo"


func HandleHelpCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
    helpEmbed := &discordgo.MessageEmbed{
        Title:       "Help: List of Commands",
        Description: "Here's a list of all commands you can use:",
        Color:       0x00ff00, // Green, change as per your preference
		Fields: []*discordgo.MessageEmbedField{
            {
                Name:  "!weather <city name>",
                Value: "Displays the current weather for the specified city.",
            },
            {
                Name:  "!translate <lang> \"<text>\"",
                Value: "Translates the given text to any language you want. Example: !translate russian \"Hello, how are you?\"",
            },
            {
                Name:  "!remind <time> <message>",
                Value: "Sets a reminder. Example: !remind 10m Take a break",
            },
            {
                Name:  "!gpt <question>",
                Value: "Asks a question to GPT-3.5-turbo model.",
            },
            {
                Name:  "!rps",
                Value: "Start or join a Rock-Paper-Scissors game.",
            },
            {
                Name:  "!poll <question> | <option1>, <option2>, ...",
                Value: "Creates a new poll with the specified question and options.",
            },
            {
                Name:  "!vote <option number>",
                Value: "Votes on the most recent poll in the channel.",
            },
            {
                Name:  "!result",
                Value: "Displays the results of the most recent poll in the channel.",
            },
            {
                Name:  "!help",
                Value: "Displays this help message.",
            },
        },
        Footer: &discordgo.MessageEmbedFooter{
            Text: "Bot developed by Raiymbek Nazymkhan. For more help, contact rnazymxan@gmail.com.",
        },
    }

    s.ChannelMessageSendEmbed(m.ChannelID, helpEmbed)
}
