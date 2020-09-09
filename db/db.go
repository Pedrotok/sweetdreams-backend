package db

import (
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/net/context"
)

// InitialConnection will create new connection to mongo db
func InitialConnection(dbName string, mongoURI string) *mongo.Database {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("Error while connecting to mongo: %v, URI: %v\n", err, mongoURI)
	}
	return client.Database(dbName)
}
