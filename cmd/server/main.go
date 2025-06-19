// File: cmd/server/main.go
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/mariopaath23/backend-jte-ticketing/internal/config"
	"github.com/mariopaath23/backend-jte-ticketing/internal/database"
	apphandlers "github.com/mariopaath23/backend-jte-ticketing/internal/handlers"
	"github.com/mariopaath23/backend-jte-ticketing/internal/middleware"
)

func main() {
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("could not load config: %v", err)
	}

	// --- DIAGNOSTIC LOGGING ---
	// Log the database connection details to help debug connection issues.
	log.Println("---------------------------------------------------------")
	log.Printf("Attempting to connect to MongoDB URI: %s", cfg.MongoURI)
	log.Printf("Targeting Database: %s", cfg.MongoDatabase)
	log.Println("---------------------------------------------------------")

	db, err := database.Connect(cfg.MongoURI, cfg.MongoDatabase)
	if err != nil {
		log.Fatalf("FATAL: Could not connect to database: %v", err)
	}
	log.Println("SUCCESS: Connection to MongoDB established.")

	// Initialize all handlers
	userHandler := apphandlers.NewUserHandler(db)
	statusHandler := apphandlers.NewStatusHandler(db)
	announcementHandler := apphandlers.NewAnnouncementHandler(db)
	catalogHandler := apphandlers.NewCatalogHandler(db)
	reservationHandler := apphandlers.NewReservationHandler(db)

	r := mux.NewRouter()
	api := r.PathPrefix("/api").Subrouter()

	// --- Public Routes ---
	api.HandleFunc("/register", userHandler.Register).Methods("POST")
	api.HandleFunc("/login", userHandler.Login).Methods("POST")
	api.HandleFunc("/logout", userHandler.Logout).Methods("POST")
	api.HandleFunc("/status/rooms", statusHandler.GetRooms).Methods("GET")
	api.HandleFunc("/status/inventory", statusHandler.GetInventoryRequests).Methods("GET")
	api.HandleFunc("/announcements", announcementHandler.GetAnnouncements).Methods("GET")
	api.HandleFunc("/catalog/search", catalogHandler.SearchCatalog).Methods("GET")
	api.HandleFunc("/catalog/room/{id}", catalogHandler.GetRoomByID).Methods("GET")

	// --- Protected Routes ---
	api.Handle("/reservations", middleware.Auth(http.HandlerFunc(reservationHandler.CreateReservation))).Methods("POST")
	api.Handle("/validate-token", middleware.Auth(http.HandlerFunc(userHandler.ValidateToken))).Methods("GET")
	api.Handle("/login-logs", middleware.Auth(http.HandlerFunc(userHandler.GetLoginLogs))).Methods("GET")

	// --- CORS Configuration ---
	allowedOrigins := handlers.AllowedOrigins([]string{"http://localhost:3000"})
	allowedMethods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})
	allowedHeaders := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization", "Credentials"})
	allowCredentials := handlers.AllowCredentials()

	port := cfg.APIPort
	fmt.Printf("Server starting on port %s...\n", port)

	log.Fatal(http.ListenAndServe(":"+port, handlers.CORS(allowedOrigins, allowedMethods, allowedHeaders, allowCredentials)(r)))
}
