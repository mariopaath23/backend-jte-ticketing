package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// LoginLog represents a single login event for a user.
type LoginLog struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	UserID    primitive.ObjectID `bson:"user_id" json:"user_id"`
	Timestamp time.Time          `bson:"timestamp" json:"timestamp"`
	UserAgent string             `bson:"user_agent" json:"user_agent"`
}
