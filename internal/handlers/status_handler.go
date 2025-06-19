package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/mariopaath23/backend-jte-ticketing/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// StatusHandler handles requests for status page data.
type StatusHandler struct {
	db *mongo.Database
}

// NewStatusHandler creates a new StatusHandler.
func NewStatusHandler(db *mongo.Database) *StatusHandler {
	return &StatusHandler{db: db}
}

// GetRooms fetches all rooms from the database.
func (h *StatusHandler) GetRooms(w http.ResponseWriter, r *http.Request) {
	roomsCollection := h.db.Collection("rooms")
	cursor, err := roomsCollection.Find(context.TODO(), bson.D{})
	if err != nil {
		http.Error(w, "Failed to retrieve rooms", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.TODO())

	var rooms []models.Room
	if err = cursor.All(context.TODO(), &rooms); err != nil {
		http.Error(w, "Failed to parse rooms data", http.StatusInternalServerError)
		return
	}

	if rooms == nil {
		rooms = []models.Room{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rooms)
}

// GetInventoryRequests fetches all inventory requests.
func (h *StatusHandler) GetInventoryRequests(w http.ResponseWriter, r *http.Request) {
	inventoryCollection := h.db.Collection("inventory_requests")

	findOptions := options.Find()
	findOptions.SetSort(bson.D{{Key: "request_date", Value: -1}}) // Sort by most recent request

	cursor, err := inventoryCollection.Find(context.TODO(), bson.D{}, findOptions)
	if err != nil {
		http.Error(w, "Failed to retrieve inventory requests", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.TODO())

	var requests []models.InventoryRequest
	if err = cursor.All(context.TODO(), &requests); err != nil {
		http.Error(w, "Failed to parse inventory requests data", http.StatusInternalServerError)
		return
	}

	if requests == nil {
		requests = []models.InventoryRequest{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(requests)
}
