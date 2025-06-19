package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/mariopaath23/backend-jte-ticketing/internal/config"
	"github.com/mariopaath23/backend-jte-ticketing/internal/database"
	"github.com/mariopaath23/backend-jte-ticketing/internal/seeds"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// Load environment variables
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("could not load config: %v", err)
	}

	// Connect to the database
	db, err := database.Connect(cfg.MongoURI, cfg.MongoDatabase)
	if err != nil {
		log.Fatalf("could not connect to database: %v", err)
	}

	// Check for command-line arguments
	if len(os.Args) < 2 {
		log.Fatal("Please provide a command: 'migrate' or 'seed'")
	}

	command := os.Args[1]

	switch command {
	case "migrate":
		fmt.Println("Running migrations...")
		migrateUsersCollection(db)
		migrateLoginLogsCollection(db)
		migrateRoomsCollection(db)
		migrateInventoryRequestsCollection(db)
		migrateAnnouncementsCollection(db)
		migrateReservationsCollection(db)
		fmt.Println("Migrations completed successfully.")
	case "seed":
		fmt.Println("Running seeders...")
		seeds.SeedUsers(db)
		seeds.SeedStatusData(db)
		seeds.SeedAnnouncements(db)
		fmt.Println("Seeding completed successfully.")
	default:
		log.Fatalf("Unknown command: %s. Available commands: 'migrate', 'seed'", command)
	}
}

// migrateUsersCollection ensures the users collection has the correct indexes.
func migrateUsersCollection(db *mongo.Database) {
	usersCollection := db.Collection("users")

	// Create a unique index on the 'email' field.
	// This prevents duplicate emails and speeds up queries by email.
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "email", Value: 1}},
		Options: options.Index().SetUnique(true),
	}

	// Create the index
	_, err := usersCollection.Indexes().CreateOne(context.TODO(), indexModel)
	if err != nil {
		log.Fatalf("Failed to create index on 'email': %v", err)
	}

	fmt.Println("Successfully created unique index on 'email' field in 'users' collection.")
}

func migrateLoginLogsCollection(db *mongo.Database) {
	// ... (no changes here)
	loginLogsCollection := db.Collection("login_logs")
	indexModel := mongo.IndexModel{
		Keys: bson.D{
			{Key: "user_id", Value: 1},
			{Key: "timestamp", Value: -1},
		},
	}
	_, err := loginLogsCollection.Indexes().CreateOne(context.TODO(), indexModel)
	if err != nil {
		log.Fatalf("Failed to create index on 'login_logs': %v", err)
	}
	fmt.Println("Successfully created index on 'user_id' and 'timestamp' fields in 'login_logs' collection.")
}

// migrateRoomsCollection creates indexes for the rooms collection.
func migrateRoomsCollection(db *mongo.Database) {
	roomsCollection := db.Collection("rooms")
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "room_id", Value: 1}},
		Options: options.Index().SetUnique(true),
	}
	_, err := roomsCollection.Indexes().CreateOne(context.TODO(), indexModel)
	if err != nil {
		log.Fatalf("Failed to create index on 'room_id': %v", err)
	}
	fmt.Println("Successfully created unique index on 'room_id' field in 'rooms' collection.")
}

// migrateInventoryRequestsCollection creates indexes for the inventory_requests collection.
func migrateInventoryRequestsCollection(db *mongo.Database) {
	inventoryRequestsCollection := db.Collection("inventory_requests")
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "request_id", Value: 1}},
		Options: options.Index().SetUnique(true),
	}
	_, err := inventoryRequestsCollection.Indexes().CreateOne(context.TODO(), indexModel)
	if err != nil {
		log.Fatalf("Failed to create index on 'request_id': %v", err)
	}
	fmt.Println("Successfully created unique index on 'request_id' field in 'inventory_requests' collection.")
}

func migrateAnnouncementsCollection(db *mongo.Database) {
	collection := db.Collection("announcements")
	// Index on date_published for fast sorting by newest
	indexModel := mongo.IndexModel{
		Keys: bson.D{{Key: "date_published", Value: -1}},
	}
	_, err := collection.Indexes().CreateOne(context.TODO(), indexModel)
	if err != nil {
		log.Fatalf("Failed to create index on 'date_published': %v", err)
	}
	fmt.Println("Successfully created index on 'date_published' field in 'announcements' collection.")
}

func migrateReservationsCollection(db *mongo.Database) {
	collection := db.Collection("reservations")

	indexModel := mongo.IndexModel{
		Keys: bson.D{
			{Key: "room_id", Value: 1},
			{Key: "start_time", Value: 1},
			{Key: "end_time", Value: 1},
		},
	}

	_, err := collection.Indexes().CreateOne(context.TODO(), indexModel)
	if err != nil {
		log.Fatalf("Failed to create index on 'reservations' collection: %v", err)
	}
	fmt.Println("Successfully created index on 'reservations' collection.")
}
