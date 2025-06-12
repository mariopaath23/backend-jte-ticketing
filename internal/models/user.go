package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// User represents a user in the database.
type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Email    string             `bson:"email" json:"email"`
	Password string             `bson:"password" json:"password"` // This will store the hashed password
}

// Credentials is used for parsing login and registration requests.
type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
