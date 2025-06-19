package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/mariopaath23/backend-jte-ticketing/internal/auth"
	"github.com/mariopaath23/backend-jte-ticketing/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AnnouncementHandler struct {
	db *mongo.Database
}

func NewAnnouncementHandler(db *mongo.Database) *AnnouncementHandler {
	return &AnnouncementHandler{db: db}
}

func (h *AnnouncementHandler) GetAnnouncements(w http.ResponseWriter, r *http.Request) {
	collection := h.db.Collection("announcements")
	filter := bson.M{"announcement_type": "public"}

	// Check for an authorization token
	authHeader := r.Header.Get("Authorization")
	if authHeader != "" {
		parts := strings.Split(authHeader, " ")
		if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
			tokenString := parts[1]

			claims, err := auth.ValidateJWT(tokenString)
			if err == nil && claims.Role == "admin" {
				filter = bson.M{}
			}
		}
	}

	// Sort by most recent
	findOptions := options.Find()
	findOptions.SetSort(bson.D{{Key: "date_published", Value: -1}})

	cursor, err := collection.Find(context.TODO(), filter, findOptions)
	if err != nil {
		http.Error(w, "Failed to retrieve announcements", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.TODO())

	var announcements []models.Announcement
	if err = cursor.All(context.TODO(), &announcements); err != nil {
		http.Error(w, "Failed to parse announcements data", http.StatusInternalServerError)
		return
	}

	if announcements == nil {
		announcements = []models.Announcement{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(announcements)
}
