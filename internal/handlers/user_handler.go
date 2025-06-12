package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/mariopaath23/backend-jte-ticketing/internal/auth"
	"github.com/mariopaath23/backend-jte-ticketing/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

// UserHandler handles user-related HTTP requests.
type UserHandler struct {
	db *mongo.Database
}

// NewUserHandler creates a new UserHandler.
func NewUserHandler(db *mongo.Database) *UserHandler {
	return &UserHandler{db: db}
}

// Register creates a new user.
func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var creds models.Credentials
	// Decode the request body into the Credentials struct
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Hash the user's password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(creds.Password), 8)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Create a new user model
	newUser := models.User{
		ID:       primitive.NewObjectID(),
		Email:    creds.Email,
		Password: string(hashedPassword),
	}

	// Insert the new user into the database
	collection := h.db.Collection("users")
	_, err = collection.InsertOne(context.TODO(), newUser)
	if err != nil {
		// Handle potential duplicate email error
		if mongo.IsDuplicateKeyError(err) {
			http.Error(w, "Email address already in use", http.StatusConflict)
			return
		}
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	// Respond with a success message
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User created successfully"})
}

// Login authenticates a user and returns a JWT.
func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var creds models.Credentials
	// Decode the request body
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var user models.User
	collection := h.db.Collection("users")
	// Find the user by email
	err := collection.FindOne(context.TODO(), bson.M{"email": creds.Email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "Invalid email or password", http.StatusUnauthorized)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Compare the stored hashed password with the password from the request
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password)); err != nil {
		// Passwords don't match
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	// If passwords match, generate a JWT token
	tokenString, err := auth.GenerateJWT(user.Email)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Set the token in a cookie or return it in the response body
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    tokenString,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true, // Makes the cookie inaccessible to client-side scripts
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}
