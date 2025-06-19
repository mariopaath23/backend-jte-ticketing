package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Announcement represents a single announcement post.
type Announcement struct {
	ID               primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Title            string             `bson:"title" json:"title"`
	Author           string             `bson:"author" json:"author"`
	DatePublished    time.Time          `bson:"date_published" json:"date_published"`
	Content          string             `bson:"content" json:"content"`
	Tags             []string           `bson:"tags,omitempty" json:"tags,omitempty"`
	AnnouncementType string             `bson:"announcement_type" json:"announcement_type"` // "public" or "private"
}
