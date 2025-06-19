package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/mariopaath23/backend-jte-ticketing/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// CatalogHandler handles requests for catalog data, including search.
type CatalogHandler struct {
	db *mongo.Database
}

// NewCatalogHandler creates a new CatalogHandler.
func NewCatalogHandler(db *mongo.Database) *CatalogHandler {
	return &CatalogHandler{db: db}
}

// GetRoomByID fetches a single room by its RoomID.
func (h *CatalogHandler) GetRoomByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	// Convert the ID string from the URL to a MongoDB ObjectID
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "Invalid Room ID format", http.StatusBadRequest)
		return
	}

	collection := h.db.Collection("rooms")
	var room models.Room

	// Find the room where the '_id' field matches the ObjectID
	err = collection.FindOne(context.TODO(), bson.M{"_id": objID}).Decode(&room)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "Room not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to retrieve room data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(room)
}

// SearchCatalog fetches rooms based on search query and status filters.
func (h *CatalogHandler) SearchCatalog(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	searchQuery := strings.TrimSpace(query.Get("q"))
	statusFilter := query.Get("status")
	typeFilter := query.Get("type")

	// For now, we only support searching for "ruangan".
	if typeFilter != "" && typeFilter != "ruangan" {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]models.Room{}) // Return empty for non-room searches
		return
	}

	collection := h.db.Collection("rooms")
	filter := bson.M{}

	// Add search query to filter (case-insensitive regex search on the 'name' field)
	if searchQuery != "" {
		filter["name"] = bson.M{"$regex": searchQuery, "$options": "i"}
	}

	// Add status filter if provided
	if statusFilter != "" {
		if statusFilter == "tersedia" {
			filter["status"] = "Available"
		} else if statusFilter == "tidak tersedia" {
			filter["status"] = bson.M{"$in": []string{"In Use", "Under Maintenance"}}
		}
	}

	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		http.Error(w, "Failed to execute search", http.StatusInternalServerError)
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
