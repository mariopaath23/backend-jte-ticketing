package seeds

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/mariopaath23/backend-jte-ticketing/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// SeedStatusData populates the database with initial rooms and inventory requests.
func SeedStatusData(db *mongo.Database) {
	fmt.Println("Seeding status data...")
	seedRooms(db)
	seedInventoryRequests(db)
	fmt.Println("Status data seeding complete.")
}

func seedRooms(db *mongo.Database) {
	roomsCollection := db.Collection("rooms")
	// Expanded and corrected list of rooms based on ruanganDummy.ts
	rooms := []models.Room{
		{RoomID: "R001", Name: "Auditorium Dekanat", ImageURL: "/assets/ruangan/Auditorium.jpg", Status: "Available", Capacity: 150, Location: "Gedung Dekanat Fakultas Teknik, Lantai 5", Type: "high", Facility: models.Facility{FurnitureAvailable: true, DisplayAvailable: true, AudioAvailable: true, ACAvailable: true}},
		{RoomID: "R002", Name: "Creative Room", ImageURL: "/assets/ruangan/Creative1.jpg", Status: "In Use", Capacity: 150, Location: "Gedung Jurusan Teknik Elektro, Lantai 2", Type: "medium", Facility: models.Facility{FurnitureAvailable: true, DisplayAvailable: true, AudioAvailable: true, ACAvailable: true}},
		{RoomID: "R003", Name: "JTE-1", ImageURL: "", Status: "In Use", Capacity: 30, Location: "Gedung Jurusan Teknik Elektro, Lantai 1", Type: "low", Facility: models.Facility{FurnitureAvailable: true, DisplayAvailable: true, AudioAvailable: false, ACAvailable: true}},
		{RoomID: "R004", Name: "JTE-2", ImageURL: "", Status: "In Use", Capacity: 30, Location: "Gedung Jurusan Teknik Elektro, Lantai 1", Type: "low", Facility: models.Facility{FurnitureAvailable: true, DisplayAvailable: true, AudioAvailable: false, ACAvailable: true}},
		{RoomID: "R005", Name: "JTE-3", ImageURL: "", Status: "In Use", Capacity: 30, Location: "Gedung Jurusan Teknik Elektro, Lantai 1", Type: "low", Facility: models.Facility{FurnitureAvailable: true, DisplayAvailable: true, AudioAvailable: false, ACAvailable: false}},
		{RoomID: "R006", Name: "JTE-4", ImageURL: "", Status: "In Use", Capacity: 30, Location: "Gedung Jurusan Teknik Elektro, Lantai 1", Type: "low", Facility: models.Facility{FurnitureAvailable: true, DisplayAvailable: true, AudioAvailable: false, ACAvailable: true}},
		{RoomID: "R007", Name: "JTE-5", ImageURL: "", Status: "In Use", Capacity: 30, Location: "Gedung Jurusan Teknik Elektro, Lantai 1", Type: "low", Facility: models.Facility{FurnitureAvailable: true, DisplayAvailable: true, AudioAvailable: false, ACAvailable: true}},
		{RoomID: "R008", Name: "JTE-6", ImageURL: "", Status: "In Use", Capacity: 30, Location: "Gedung Jurusan Teknik Elektro, Lantai 1", Type: "low", Facility: models.Facility{FurnitureAvailable: true, DisplayAvailable: true, AudioAvailable: false, ACAvailable: true}},
		{RoomID: "R009", Name: "JTE-7", ImageURL: "", Status: "In Use", Capacity: 30, Location: "Gedung Jurusan Teknik Elektro, Lantai 1", Type: "low", Facility: models.Facility{FurnitureAvailable: true, DisplayAvailable: true, AudioAvailable: false, ACAvailable: true}},
		{RoomID: "R010", Name: "JTE-8", ImageURL: "", Status: "In Use", Capacity: 50, Location: "Gedung Jurusan Teknik Elektro, Lantai 3", Type: "low", Facility: models.Facility{FurnitureAvailable: true, DisplayAvailable: true, AudioAvailable: true, ACAvailable: true}},
		{RoomID: "R011", Name: "JTE-9", ImageURL: "", Status: "In Use", Capacity: 50, Location: "Gedung Jurusan Teknik Elektro, Lantai 3", Type: "low", Facility: models.Facility{FurnitureAvailable: true, DisplayAvailable: true, AudioAvailable: true, ACAvailable: true}},
	}

	for i, room := range rooms {
		var existingRoom models.Room
		// Generate a consistent RoomID based on index
		room.RoomID = fmt.Sprintf("R%03d", i+1)

		// Check if a room with the same RoomID already exists
		err := roomsCollection.FindOne(context.TODO(), bson.M{"room_id": room.RoomID}).Decode(&existingRoom)
		if err == mongo.ErrNoDocuments {
			// If it doesn't exist, insert it
			room.ID = primitive.NewObjectID()
			_, insertErr := roomsCollection.InsertOne(context.TODO(), room)
			if insertErr != nil {
				log.Printf("Failed to seed room %s: %v", room.Name, insertErr)
			} else {
				fmt.Printf("Successfully seeded room: %s\n", room.Name)
			}
		} else if err != nil {
			log.Printf("Error checking for room %s: %v", room.Name, err)
		} else {
			fmt.Printf("Room %s (%s) already exists. Skipping.\n", existingRoom.Name, existingRoom.RoomID)
		}
	}
}

func seedInventoryRequests(db *mongo.Database) {
	inventoryCollection := db.Collection("inventory_requests")
	requests := []models.InventoryRequest{
		{RequestID: "REQ001", RequesterName: "John Doe", ItemName: "Laptop", RequestDate: parseDate("2024-05-01"), Status: "Approved", PickupDate: parseDate("2024-05-03")},
		{RequestID: "REQ002", RequesterName: "Jane Smith", ItemName: "Proyektor", RequestDate: parseDate("2024-05-02"), Status: "Pending", PickupDate: time.Time{}},
		{RequestID: "REQ003", RequesterName: "Peter Jones", ItemName: "Kabel HDMI", RequestDate: parseDate("2024-05-03"), Status: "Rejected", PickupDate: time.Time{}},
		{RequestID: "REQ004", RequesterName: "Mary Jane", ItemName: "Papan Tulis Spidol", RequestDate: parseDate("2024-05-04"), Status: "Approved", PickupDate: parseDate("2024-05-06")},
	}

	for _, req := range requests {
		var existingReq models.InventoryRequest
		err := inventoryCollection.FindOne(context.TODO(), bson.M{"request_id": req.RequestID}).Decode(&existingReq)
		if err == mongo.ErrNoDocuments {
			req.ID = primitive.NewObjectID()
			_, insertErr := inventoryCollection.InsertOne(context.TODO(), req)
			if insertErr != nil {
				log.Printf("Failed to seed inventory request %s: %v", req.RequestID, insertErr)
			} else {
				fmt.Printf("Successfully seeded inventory request: %s\n", req.RequestID)
			}
		}
	}
}

func parseDate(dateStr string) time.Time {
	t, _ := time.Parse("2006-01-02", dateStr)
	return t
}
