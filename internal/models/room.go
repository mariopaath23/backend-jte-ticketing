package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Facility struct {
	FurnitureAvailable bool `bson:"furniture_available" json:"furnitureAvailable"`
	DisplayAvailable   bool `bson:"display_available" json:"displayAvailable"`
	AudioAvailable     bool `bson:"audio_available" json:"audioAvailable"`
	ACAvailable        bool `bson:"ac_available" json:"acAvailable"`
}

type Room struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	RoomID   string             `bson:"room_id" json:"room_id"`
	Name     string             `bson:"name" json:"name"`
	ImageURL string             `bson:"image_url" json:"imageUrl"`
	Status   string             `bson:"status" json:"status"`
	Capacity int                `bson:"capacity" json:"capacity"`
	Location string             `bson:"location" json:"location"`
	Type     string             `bson:"type" json:"type"`
	Facility Facility           `bson:"facility,omitempty" json:"facility,omitempty"`
}
