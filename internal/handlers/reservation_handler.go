package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/mariopaath23/backend-jte-ticketing/internal/auth"
	"github.com/mariopaath23/backend-jte-ticketing/internal/middleware"
	"github.com/mariopaath23/backend-jte-ticketing/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// ReservationHandler handles requests for reservation data.
type ReservationHandler struct {
	db *mongo.Database
}

// NewReservationHandler creates a new ReservationHandler.
func NewReservationHandler(db *mongo.Database) *ReservationHandler {
	return &ReservationHandler{db: db}
}

// CreateReservation handles the creation of a new room reservation.
func (h *ReservationHandler) CreateReservation(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middleware.ClaimsKey).(*auth.Claims)
	if !ok || claims == nil {
		http.Error(w, "Unauthorized: Could not retrieve user claims.", http.StatusUnauthorized)
		return
	}

	var payload models.CreateReservationPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// --- Validation ---
	roomObjID, err := primitive.ObjectIDFromHex(payload.RoomID)
	if err != nil {
		http.Error(w, "Invalid Room ID format", http.StatusBadRequest)
		return
	}
	startTime, err := time.Parse(time.RFC3339, payload.StartTime)
	if err != nil {
		http.Error(w, "Invalid Start Time format", http.StatusBadRequest)
		return
	}
	endTime, err := time.Parse(time.RFC3339, payload.EndTime)
	if err != nil {
		http.Error(w, "Invalid End Time format", http.StatusBadRequest)
		return
	}
	if endTime.Before(startTime) || endTime.Equal(startTime) {
		http.Error(w, "End time must be after start time", http.StatusBadRequest)
		return
	}

	// --- Conflict Check ---
	collection := h.db.Collection("reservations")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{
		"room_id": roomObjID,
		"status":  "Approved",
		"$or": []bson.M{
			{"start_time": bson.M{"$gte": startTime, "$lt": endTime}},
			{"end_time": bson.M{"$gt": startTime, "$lte": endTime}},
			{"start_time": bson.M{"$lte": startTime}, "end_time": bson.M{"$gte": endTime}},
		},
	}

	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		http.Error(w, "Failed to check for booking conflicts", http.StatusInternalServerError)
		return
	}
	if count > 0 {
		http.Error(w, "The selected time slot is unavailable due to a conflict.", http.StatusConflict)
		return
	}

	// --- Create Reservation ---
	newReservation := models.Reservation{
		ID:          primitive.NewObjectID(),
		RoomID:      roomObjID,
		UserID:      claims.UserID,
		Purpose:     payload.Purpose,
		Description: payload.Description,
		StartTime:   startTime,
		EndTime:     endTime,
		Status:      "Pending",
		CreatedAt:   time.Now(),
	}

	// Use a new context with a timeout for the insert operation.
	// This makes the write operation more robust.
	insertCtx, insertCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer insertCancel()

	insertResult, err := collection.InsertOne(insertCtx, newReservation)
	if err != nil {
		log.Printf("ERROR: Failed to insert reservation into database: %v", err)
		http.Error(w, "Failed to create reservation", http.StatusInternalServerError)
		return
	}

	log.Printf("SUCCESS: Reservation created successfully with ID: %v", insertResult.InsertedID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message":       "Reservasi berhasil dibuat",
		"reservationId": newReservation.ID.Hex(),
	})
}
