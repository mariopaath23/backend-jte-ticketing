// File: internal/handlers/user_handler.go
package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/mariopaath23/backend-jte-ticketing/internal/auth"
	"github.com/mariopaath23/backend-jte-ticketing/internal/middleware" // Import the middleware package
	"github.com/mariopaath23/backend-jte-ticketing/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func (h *UserHandler) ValidateToken(w http.ResponseWriter, r *http.Request) {
	// Get user claims from context using the exported key from the middleware package.
	claims, ok := r.Context().Value(middleware.ClaimsKey).(*auth.Claims)
	if !ok {
		http.Error(w, "Could not retrieve user claims", http.StatusInternalServerError)
		return
	}

	var user models.User
	collection := h.db.Collection("users")
	err := collection.FindOne(context.TODO(), bson.M{"_id": claims.UserID}).Decode(&user)
	if err != nil {
		http.Error(w, "User not found or invalid", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  http.StatusOK,
		"message": "Token is valid",
		"user": map[string]string{
			"email": user.Email,
			"role":  user.Role,
		},
	})
}

// ... (Register, Login, Logout, and other functions remain the same) ...
// Login authenticates a user and returns a JWT and user data.
func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var creds models.Credentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var user models.User
	collection := h.db.Collection("users")
	err := collection.FindOne(context.TODO(), bson.M{"email": creds.Email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"status":  http.StatusUnauthorized,
				"message": "Email atau password salah!",
			})
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password)); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  http.StatusUnauthorized,
			"message": "Email atau password salah!",
		})
		return
	}

	tokenString, err := auth.GenerateJWT(user.ID, user.Email, user.Role)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	go h.logLoginSession(user.ID, r.UserAgent())

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    tokenString,
		Expires:  time.Now().Add(120 * time.Minute),
		HttpOnly: true,
	})

	// Return token, user data, status, and message in the response body
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  http.StatusOK,
		"message": "Berhasil Masuk!",
		"token":   tokenString,
		"user": map[string]string{
			"email": user.Email,
			"role":  user.Role,
		},
	})
}

// Register creates a new user.
func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var creds models.Credentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(creds.Password), 8)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	newUser := models.User{
		ID:       primitive.NewObjectID(),
		Email:    creds.Email,
		Password: string(hashedPassword),
		Role:     "student",
	}

	collection := h.db.Collection("users")
	_, err = collection.InsertOne(context.TODO(), newUser)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			http.Error(w, "Email address already in use", http.StatusConflict)
			return
		}
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User created successfully"})
}

// Logout handles user logout by clearing the authentication cookie.
func (h *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		MaxAge:   -1,
	})

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Successfully logged out"})
}

func (h *UserHandler) logLoginSession(userID primitive.ObjectID, userAgent string) {
	logCollection := h.db.Collection("login_logs")

	logEntry := models.LoginLog{
		ID:        primitive.NewObjectID(),
		UserID:    userID,
		Timestamp: time.Now(),
		UserAgent: userAgent,
	}

	_, err := logCollection.InsertOne(context.TODO(), logEntry)
	if err != nil {
		log.Printf("Failed to create login log for user %s: %v", userID.Hex(), err)
	}
}

func (h *UserHandler) GetLoginLogs(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middleware.ClaimsKey).(*auth.Claims)
	if !ok {
		http.Error(w, "Could not retrieve user claims", http.StatusInternalServerError)
		return
	}

	logCollection := h.db.Collection("login_logs")

	findOptions := options.Find()
	findOptions.SetSort(bson.D{{Key: "timestamp", Value: -1}})

	cursor, err := logCollection.Find(context.TODO(), bson.M{"user_id": claims.UserID}, findOptions)
	if err != nil {
		http.Error(w, "Failed to retrieve login logs", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.TODO())

	var logs []models.LoginLog
	if err = cursor.All(context.TODO(), &logs); err != nil {
		http.Error(w, "Failed to parse login logs", http.StatusInternalServerError)
		return
	}

	if logs == nil {
		logs = []models.LoginLog{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(logs)
}
