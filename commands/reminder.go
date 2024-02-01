package commands

import (
	"context"
	db "discord-bot/database"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Reminder struct {
    ID        string    `bson:"_id,omitempty"`
    UserID    string    `bson:"user_id"`
    ChannelID string    `bson:"channel_id"`
    Message   string    `bson:"message"`
    Time      time.Time `bson:"time"`
    Reminded  bool      `bson:"reminded"`
}

func HandleRemindCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
    // Parse command
    content := strings.TrimPrefix(m.Content, "!remind ")
    parts := strings.SplitN(content, " ", 2)
    if len(parts) < 2 {
        s.ChannelMessageSend(m.ChannelID, "Usage: !remind [time] [message]")
        return
    }

    duration, err := time.ParseDuration(parts[0])
    if err != nil {
        s.ChannelMessageSend(m.ChannelID, "Invalid time format. Please use '10m', '2h45m', etc.")
        return
    }

    reminderTime := time.Now().Add(duration)
	reminder := Reminder{
		UserID:    m.Author.ID,
		ChannelID: m.ChannelID,
		Message:   parts[1],
		Time:      reminderTime,
		Reminded:  false, 
	}

    // Store reminder in MongoDB
    collection := db.MongoClient.Database("DISCORD-BOT").Collection("reminders")
    _, err = collection.InsertOne(context.Background(), reminder)
    if err != nil {
        log.Println("Failed to insert reminder:", err)
        s.ChannelMessageSend(m.ChannelID, "Failed to set reminder.")
        return
    }

    s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("I will remind you in %s to %s.", parts[0], parts[1]))
}

func RemindScheduler(s *discordgo.Session) {
    ticker := time.NewTicker(5 * time.Second)
    for range ticker.C {
        now := time.Now()
        collection := db.MongoClient.Database("DISCORD-BOT").Collection("reminders")
        cursor, err := collection.Find(context.Background(), bson.M{"time": bson.M{"$lte": now}, "reminded": false})

        if err != nil {
            log.Println("Error fetching due reminders:", err)
            continue
        }

        var reminders []Reminder
        if err = cursor.All(context.Background(), &reminders); err != nil {
            log.Println("Error decoding reminders:", err)
            continue
        }

		for _, reminder := range reminders {
			log.Printf("Sending reminder to user %s in channel %s\n", reminder.UserID, reminder.ChannelID)
			s.ChannelMessageSend(reminder.ChannelID, fmt.Sprintf("<@%s>, remember to %s!", reminder.UserID, reminder.Message))
			log.Println("Reminder sent successfully")
			objectID, err := primitive.ObjectIDFromHex(reminder.ID) 
			if err != nil {
				log.Println("Error converting ID to ObjectID:", err)
				continue
			}
		
			updateResult, err := collection.UpdateOne(
				context.Background(),
				bson.M{"_id": objectID},
				bson.M{"$set": bson.M{"reminded": true}},
			)
			if err != nil {
				log.Println("Failed to update reminder:", err)
			} else if updateResult.ModifiedCount == 0 {
				log.Println("No reminder was updated, check the ID:", reminder.ID)
			} else {
				log.Printf("Reminder updated successfully, ID: %s\n", reminder.ID)
			}
		}
    }
}

