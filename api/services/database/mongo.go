package database

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoService struct {
	connection *mongo.Client
}

func NewMongoService(uri string) (*MongoService, error) {
	var err error
	clientOptions := options.Client().ApplyURI(uri)

	// Connect to MongoDB
	ctx := context.TODO()
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	log.Println("Connected to MongoDB!")
	return &MongoService{connection: client}, nil
}

// GetCollection returns a MongoDB collection
func (mc MongoService) GetCollection(databaseName, collectionName string) *mongo.Collection {
	return mc.connection.Database(databaseName).Collection(collectionName)
}

// Disconnect closes the connection to the MongoDB server
func (mc MongoService) Disconnect() error {
	if err := mc.connection.Disconnect(context.TODO()); err != nil {
		return fmt.Errorf("failed to disconnect from MongoDB: %w", err)
	}
	log.Println("Disconnected from MongoDB!")
	return nil
}

func (ms *MongoService) GetClient() *mongo.Client {
	if ms.connection == nil {
		log.Fatal("MongoDB connection is nil!")
	}
	return ms.connection
}
