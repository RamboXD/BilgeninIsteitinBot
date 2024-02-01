package commands

import (
	"context"
	db "discord-bot/database"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Poll struct {
    ID        primitive.ObjectID `bson:"_id,omitempty"`
    Question  string             `bson:"question"`
    Options   []string           `bson:"options"`
    Votes     []Vote             `bson:"votes"`
    ChannelID string             `bson:"channel_id"` 
}

type Vote struct {
    UserID string `bson:"user_id"`
    Option int    `bson:"option"`
}

func HandlePollCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
    content := strings.TrimPrefix(m.Content, "!poll ")
    parts := strings.SplitN(content, "|", 2)
    if len(parts) < 2 {
        s.ChannelMessageSend(m.ChannelID, "Usage: !poll Question | option1, option2, option3...")
        return
    }

    question := parts[0]
    options := strings.Split(parts[1], ",")
    if len(options) < 2 {
        s.ChannelMessageSend(m.ChannelID, "You must provide at least two options.")
        return
    }

	poll := Poll{
		Question:  question,
		Options:   options,
		Votes:     []Vote{},
		ChannelID: m.ChannelID,
	}
	

    collection := db.MongoClient.Database("DISCORD-BOT").Collection("polls")
    _, err := collection.InsertOne(context.Background(), poll)
    if err != nil {
        // log.Println("Failed to insert poll:", err)
        s.ChannelMessageSend(m.ChannelID, "Failed to create poll.")
        return
    }

    s.ChannelMessageSend(m.ChannelID, "Poll created successfully.")
}

func HandleVoteCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
    // Extract the vote option from the command
    voteOption, err := strconv.Atoi(strings.TrimSpace(strings.TrimPrefix(m.Content, "!vote")))
    if err != nil || voteOption < 1 {
        s.ChannelMessageSend(m.ChannelID, "Please provide a valid option number to vote.")
        return
    }

    collection := db.MongoClient.Database("DISCORD-BOT").Collection("polls")
    // Find the most recent poll in the channel
    var poll Poll
	findOptions := options.FindOne().SetSort(bson.D{primitive.E{Key: "_id", Value: -1}})
    err = collection.FindOne(context.Background(), bson.M{"channel_id": m.ChannelID}, findOptions).Decode(&poll)
    if err != nil {
        s.ChannelMessageSend(m.ChannelID, "Failed to find an active poll in this channel.")
        return
    }

    // Check if the user has already voted
    for _, vote := range poll.Votes {
        if vote.UserID == m.Author.ID {
            s.ChannelMessageSend(m.ChannelID, "You have already voted in this poll.")
            return
        }
    }

    // Check if the vote option is valid
    if voteOption > len(poll.Options) {
        s.ChannelMessageSend(m.ChannelID, "Invalid vote option. Please choose a listed option.")
        return
    }

    // Add the vote to the poll
    vote := Vote{
        UserID: m.Author.ID,
        Option: voteOption - 1, // Adjust for zero-based index
    }
    update := bson.M{"$push": bson.M{"votes": vote}}
    _, err = collection.UpdateOne(context.Background(), bson.M{"_id": poll.ID}, update)
	if err != nil {
		log.Printf("Failed to record your vote: %v", err)
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Failed to record your vote: %v", err))
		return
	}

    s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Your vote for '%s' has been recorded.", poll.Options[voteOption-1]))
}

func HandleResultCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
    collection := db.MongoClient.Database("DISCORD-BOT").Collection("polls")

    // Find the most recent poll in the channel
    var poll Poll
    findOptions := options.FindOne().SetSort(bson.D{{Key: "_id", Value: -1}})
    err := collection.FindOne(context.Background(), bson.M{"channel_id": m.ChannelID}, findOptions).Decode(&poll)
    if err != nil {
        s.ChannelMessageSend(m.ChannelID, "Failed to find an active poll in this channel.")
        return
    }

    // Tally the votes
    voteCount := make(map[string]int)
    for _, vote := range poll.Votes {
        option := poll.Options[vote.Option]
        voteCount[option]++
    }

    // Format the results
    var results strings.Builder
    results.WriteString(fmt.Sprintf("Results for the poll: %s\n", poll.Question))
    for i, option := range poll.Options {
        count := voteCount[option]
        results.WriteString(fmt.Sprintf("%d. %s - %d votes\n", i+1, option, count))
    }

    // Send the results to the channel
    s.ChannelMessageSend(m.ChannelID, results.String())
}