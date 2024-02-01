package db

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Store Bot API Tokens:
var (
	MongoDBURI       string
)

var MongoClient *mongo.Client

func InitMongoDB() {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    var err error
    MongoClient, err = mongo.Connect(ctx, options.Client().ApplyURI(MongoDBURI))
    if err != nil {
        log.Fatalf("Failed to connect to MongoDB: %v", err)
    }
    // Check the connection
    err = MongoClient.Ping(ctx, nil)
    if err != nil {
        log.Fatalf("Failed to ping MongoDB: %v", err)
    }
    log.Println("Connected to MongoDB successfully")
}
