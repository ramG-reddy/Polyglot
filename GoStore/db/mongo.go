package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	// Collection name in MongoDB
	SMSRecordsCollection = "sms_records"
)

var (
	// Client is the MongoDB client instance
	Client *mongo.Client
	// Database is the SMS Store database
	Database *mongo.Database
)

// InitMongoDB establishes connection to MongoDB with retry logic
func InitMongoDB(uri, dbName string) error {
	log.Println("Initializing MongoDB connection...")

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	// Set client options
	clientOptions := options.Client().ApplyURI(uri).
		SetMaxPoolSize(50).
		SetMinPoolSize(10).
		SetMaxConnIdleTime(30 * time.Second).
		SetServerSelectionTimeout(10 * time.Second)

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Ping the database to verify connection
	pingCtx, pingCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer pingCancel()

	if err := client.Ping(pingCtx, nil); err != nil {
		return fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	Client = client
	Database = client.Database(dbName)

	log.Printf("Successfully connected to MongoDB database: %s", dbName)
	return nil
}

// CreateIndexes creates necessary indexes on the sms_records collection
func CreateIndexes() error {
	log.Println("Creating MongoDB indexes...")

	collection := Database.Collection(SMSRecordsCollection)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// First, drop all existing indexes except _id
	indexView := collection.Indexes()
	cursor, err := indexView.List(ctx)
	if err != nil {
		log.Printf("Warning: Failed to list indexes: %v", err)
	} else {
		var existingIndexes []bson.M
		if err = cursor.All(ctx, &existingIndexes); err != nil {
			log.Printf("Warning: Failed to decode indexes: %v", err)
		} else {
			for _, idx := range existingIndexes {
				indexName := idx["name"].(string)
				// Don't drop the default _id index
				if indexName != "_id_" {
					log.Printf("Dropping existing index: %s", indexName)
					if _, err := indexView.DropOne(ctx, indexName); err != nil {
						log.Printf("Warning: Failed to drop index %s: %v", indexName, err)
					}
				}
			}
		}
	}

	// Define indexes
	indexes := []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "user_id", Value: 1}},
			Options: options.Index().
				SetName("idx_user_id"),
		},
		{
			Keys: bson.D{{Key: "created_at", Value: -1}},
			Options: options.Index().
				SetName("idx_created_at"),
		},
		{
			Keys: bson.D{
				{Key: "user_id", Value: 1},
				{Key: "created_at", Value: -1},
			},
			Options: options.Index().
				SetName("idx_user_id_created_at"),
		},
	}

	// Create indexes
	indexNames, err := collection.Indexes().CreateMany(ctx, indexes)
	if err != nil {
		return fmt.Errorf("failed to create indexes: %w", err)
	}

	log.Printf("Successfully created %d indexes: %v", len(indexNames), indexNames)
	return nil
}

// GetCollection returns the sms_records collection
func GetCollection() *mongo.Collection {
	return Database.Collection(SMSRecordsCollection)
}

// Close closes the MongoDB connection gracefully
func Close() error {
	if Client == nil {
		return nil
	}

	log.Println("Closing MongoDB connection...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := Client.Disconnect(ctx); err != nil {
		return fmt.Errorf("failed to disconnect from MongoDB: %w", err)
	}

	log.Println("MongoDB connection closed successfully")
	return nil
}

// HealthCheck verifies MongoDB connection is alive
func HealthCheck() error {
	if Client == nil {
		return fmt.Errorf("MongoDB client is not initialized")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := Client.Ping(ctx, nil); err != nil {
		return fmt.Errorf("MongoDB health check failed: %w", err)
	}

	return nil
}
