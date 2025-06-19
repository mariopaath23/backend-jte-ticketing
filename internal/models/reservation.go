package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Reservation struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	RoomID      primitive.ObjectID `bson:"room_id" json:"roomId"`
	UserID      primitive.ObjectID `bson:"user_id" json:"userId"`
	Purpose     string             `bson:"purpose" json:"purpose"`
	Description string             `bson:"description,omitempty" json:"description,omitempty"`
	StartTime   time.Time          `bson:"start_time" json:"startTime"`
	EndTime     time.Time          `bson:"end_time" json:"endTime"`
	Status      string             `bson:"status" json:"status"`
	CreatedAt   time.Time          `bson:"created_at" json:"createdAt"`
}

type CreateReservationPayload struct {
	RoomID      string `json:"roomId"`
	Purpose     string `json:"purpose"`
	Description string `json:"description"`
	StartTime   string `json:"startTime"`
	EndTime     string `json:"endTime"`
}
