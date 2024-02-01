package commands

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/bwmarrin/discordgo"
)
var (
    rpsPlayer1 *discordgo.User
    rpsPlayer2 *discordgo.User
    rpsChannelID string 
)

func HandleRPSCommand(s *discordgo.Session, m *discordgo.MessageCreate) {

    // Player agreement and game start logic
    if rpsPlayer1 == nil {
        rpsPlayer1 = m.Author
        rpsChannelID = m.ChannelID

        s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s wants to play Rock-Paper-Scissors! Type `!rps` to play against them.", m.Author.Mention()))
        return
    }

    if rpsPlayer2 == nil && m.Author.ID != rpsPlayer1.ID {
        rpsPlayer2 = m.Author
        // Both players are ready, play the game
        playRPSGame(s)
        // Reset for the next game
        rpsPlayer1, rpsPlayer2, rpsChannelID = nil, nil, ""
        return
    }
}

func playRPSGame(s *discordgo.Session) {
    choices := []string{"rock", "paper", "scissors"}
    rand.Seed(time.Now().UnixNano())
    choice1 := choices[rand.Intn(len(choices))]
    choice2 := choices[rand.Intn(len(choices))]

    result := fmt.Sprintf("%s got %s, %s got %s. ", rpsPlayer1.Mention(), choice1, rpsPlayer2.Mention(), choice2)
    switch {
    case choice1 == choice2:
        result += "It's a tie!"
    case (choice1 == "rock" && choice2 == "scissors") || (choice1 == "scissors" && choice2 == "paper") || (choice1 == "paper" && choice2 == "rock"):
        result += fmt.Sprintf("%s wins!", rpsPlayer1.Mention())
    default:
        result += fmt.Sprintf("%s wins!", rpsPlayer2.Mention())
    }

    s.ChannelMessageSend(rpsChannelID, result)
}