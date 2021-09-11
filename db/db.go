package db

import (
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
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

// SetIndexes will create mongo indexes for collection and keys that sent.
func SetIndexes(collection *mongo.Collection, keys bsonx.Doc) {
	index := mongo.IndexModel{}
	index.Keys = keys
	unique := true
	index.Options = &options.IndexOptions{
		Unique: &unique,
	}
	opts := options.CreateIndexes().SetMaxTime(10 * time.Second)
	_, err := collection.Indexes().CreateOne(context.Background(), index, opts)
	if err != nil {
		log.Fatalf("Error while creating indexes: %v", err)
	}
}
