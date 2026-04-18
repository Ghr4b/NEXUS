package main

import (
	"context"
	"fmt"
	"log"
	"time"
)

// InitDB purges the users collection and seeds it with 20 normal users (role 1) and 1 admin (role 0).
func InitDB() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	collection := DB.Collection("users")

	// Purge the entire users collection
	err := collection.Drop(ctx)
	if err != nil {
		log.Fatalf("Failed to purge users collection: %v", err)
	}
	fmt.Println("Users collection purged.")

	// Seed 20 normal users (role 1)
	var users []interface{}
	for i := 1; i <= 20; i++ {
		users = append(users, User{
			Username:  fmt.Sprintf("user%d", i),
			Password:  fmt.Sprintf("password%04d", i),
			Role:      1,
			CreatedAt: time.Now(),
		})
	}

	// Seed 1 admin user (role 0)
	users = append(users, User{
		Username:  "ghassenboulares",
		Password:  "babajabliballoun123!!",
		Role:      0,
		CreatedAt: time.Now(),
	})

	_, err = collection.InsertMany(ctx, users)
	if err != nil {
		log.Fatalf("Failed to seed users: %v", err)
	}

	fmt.Println("Database seeded: 20 normal users + 1 admin.")
}
