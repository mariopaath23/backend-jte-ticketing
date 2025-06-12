package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mariopaath23/backend-jte-ticketing/internal/config"
	"github.com/mariopaath23/backend-jte-ticketing/internal/database"
	"github.com/mariopaath23/backend-jte-ticketing/internal/handlers"
	"github.com/mariopaath23/backend-jte-ticketing/internal/middleware"
)

func main() {
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("config loading failed: %v", err)
	}

	// Connect to MongoDB
	db, err := database.Connect(cfg.MongoURI, cfg.MongoDatabase)
	if err != nil {
		log.Fatalf("could not connect to database: %v", err)
	}
	// The user handler needs access to the database
	userHandler := handlers.NewUserHandler(db)

	// Create a new router
	r := mux.NewRouter()

	// Create a subrouter for API endpoints
	api := r.PathPrefix("/api").Subrouter()

	// Public routes (no authentication required)
	api.HandleFunc("/register", userHandler.Register).Methods("POST")
	api.HandleFunc("/login", userHandler.Login).Methods("POST")

	// Protected route - requires JWT authentication
	// We wrap the handler with our authentication middleware
	api.Handle("/protected", middleware.Auth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "This is a protected route! You are authenticated."}`))
	}))).Methods("GET")

	// Start the server
	port := cfg.APIPort
	fmt.Printf("Server starting on port %s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
