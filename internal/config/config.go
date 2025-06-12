package config

import (
	"github.com/joho/godotenv"
)

// Config struct holds all the configuration for the application
// The values are read by godotenv from a .env file
type Config struct {
	MongoURI      string
	MongoDatabase string
	JWTSecretKey  string
	APIPort       string
}

// LoadConfig reads configuration from file or environment variables.
func LoadConfig(path string) (config Config, err error) {
	// This will load the .env file at the given path
	// and set the environment variables
	vars, err := godotenv.Read(path + "/.env")
	if err != nil {
		return Config{}, err
	}

	config = Config{
		MongoURI:      vars["MONGO_URI"],
		MongoDatabase: vars["MONGO_DATABASE"],
		JWTSecretKey:  vars["JWT_SECRET_KEY"],
		APIPort:       vars["API_PORT"],
	}

	return
}
