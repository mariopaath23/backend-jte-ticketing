package database

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Connect establishes a connection to MongoDB and returns a database object.
func Connect(uri, dbName string) (*mongo.Database, error) {
	// Set a timeout for the connection context.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Set client options from the URI.
	clientOptions := options.Client().ApplyURI(uri)

	// Connect to MongoDB.
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	// Ping the primary node to verify that the connection is alive.
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	// Return a handle to the specified database.
	return client.Database(dbName), nil
}
