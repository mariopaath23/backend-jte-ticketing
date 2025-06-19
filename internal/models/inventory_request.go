package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// InventoryRequest represents a request for an inventory item.
// This will be used for the table on the /status page.
type InventoryRequest struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	RequestID     string             `bson:"request_id" json:"request_id"`
	RequesterName string             `bson:"requester_name" json:"requester_name"`
	ItemName      string             `bson:"item_name" json:"item_name"`
	RequestDate   time.Time          `bson:"request_date" json:"request_date"`
	Status        string             `bson:"status" json:"status"` // e.g., "Approved", "Pending", "Rejected"
	PickupDate    time.Time          `bson:"pickup_date" json:"pickup_date"`
}
