// File: internal/seeds/user_seeder.go

package seeds

import (
	"context"
	"fmt"
	"log"

	"github.com/mariopaath23/backend-jte-ticketing/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

// SeedUsers creates initial admin and student users in the database.
func SeedUsers(db *mongo.Database) {
	usersCollection := db.Collection("users")

	// --- Seed Admin User ---
	seedUser(usersCollection, "admin@jte.com", "admin123", "superadmin")

	// --- Seed Student Users ---
	seedUser(usersCollection, "user1@student.unsrat.ac.id", "password", "student")
	seedUser(usersCollection, "user2@unsrat.ac.id", "password", "admin")
}

// seedUser is a helper function to create a user if they don't already exist.
func seedUser(collection *mongo.Collection, email, password, role string) {
	// Check if the user already exists
	var existingUser models.User
	err := collection.FindOne(context.TODO(), bson.M{"email": email}).Decode(&existingUser)
	if err == nil {
		// User already exists, so we do nothing.
		fmt.Printf("User with email %s already exists. Skipping.\n", email)
		return
	}
	// If the error is not "no documents in result", it's some other problem.
	if err != mongo.ErrNoDocuments {
		log.Fatalf("Failed to check for existing user %s: %v", email, err)
	}

	// User does not exist, so we create them.
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 8)
	if err != nil {
		log.Fatalf("Failed to hash password for user %s: %v", email, err)
	}

	newUser := models.User{
		ID:       primitive.NewObjectID(),
		Email:    email,
		Password: string(hashedPassword),
		Role:     role,
	}

	_, err = collection.InsertOne(context.TODO(), newUser)
	if err != nil {
		log.Fatalf("Failed to seed user with email %s: %v", email, err)
	}

	fmt.Printf("Successfully seeded user: %s (role: %s)\n", email, role)
}
