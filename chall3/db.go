package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DB is the exported database instance
var DB *mongo.Database

// MongoClient is the exported client instance
var MongoClient *mongo.Client

// ConnectDB establishes the connection to MongoDB
func ConnectDB() {
	// Full URI with credentials
	uri := "mongodb+srv://Ghr4b@mongodb.net/?appName=Cluster0"

	clientOptions := options.Client().ApplyURI(uri)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// Ping the database to verify the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}

	fmt.Println("Successfully connected to MongoDB!")

	MongoClient = client

	// Default database we will use. Change "webapp" to your intended database name.
	DB = client.Database("webapp")
}
