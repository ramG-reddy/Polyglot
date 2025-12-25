package services

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/ramG-reddy/sms-store/db"
	"github.com/ramG-reddy/sms-store/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// SMSService handles business logic for SMS record operations
type SMSService struct {
	collection string
}

// NewSMSService creates a new SMS service instance
func NewSMSService() *SMSService {
	return &SMSService{
		collection: db.SMSRecordsCollection,
	}
}

// SaveMessage persists an SMS record to MongoDB
func (s *SMSService) SaveMessage(ctx context.Context, record *models.SMSRecord) error {
	log.Printf("Saving SMS record for user: %s", record.UserID)

	collection := db.GetCollection()

	// Set timeout for insert operation
	insertCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Insert the record
	result, err := collection.InsertOne(insertCtx, record)
	if err != nil {
		return fmt.Errorf("failed to insert SMS record: %w", err)
	}

	log.Printf("Successfully saved SMS record with ID: %v for user: %s", result.InsertedID, record.UserID)
	return nil
}

// GetMessagesByUserID retrieves all SMS messages for a specific user
// Results are sorted by created_at in descending order (newest first)
func (s *SMSService) GetMessagesByUserID(ctx context.Context, userID string) ([]*models.SMSRecord, error) {
	log.Printf("Retrieving messages for user: %s", userID)

	collection := db.GetCollection()

	// Set timeout for query operation
	queryCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// Build query filter
	filter := bson.M{"user_id": userID}

	// Set options: sort by created_at descending
	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})

	// Execute query
	cursor, err := collection.Find(queryCtx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to query messages: %w", err)
	}
	defer cursor.Close(queryCtx)

	// Decode results
	var records []*models.SMSRecord
	if err := cursor.All(queryCtx, &records); err != nil {
		return nil, fmt.Errorf("failed to decode messages: %w", err)
	}

	log.Printf("Retrieved %d messages for user: %s", len(records), userID)
	return records, nil
}

// GetRecentMessages retrieves the most recent N messages for a user
func (s *SMSService) GetRecentMessages(ctx context.Context, userID string, limit int64) ([]*models.SMSRecord, error) {
	log.Printf("Retrieving recent %d messages for user: %s", limit, userID)

	collection := db.GetCollection()

	queryCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	filter := bson.M{"user_id": userID}
	opts := options.Find().
		SetSort(bson.D{{Key: "created_at", Value: -1}}).
		SetLimit(limit)

	cursor, err := collection.Find(queryCtx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to query recent messages: %w", err)
	}
	defer cursor.Close(queryCtx)

	var records []*models.SMSRecord
	if err := cursor.All(queryCtx, &records); err != nil {
		return nil, fmt.Errorf("failed to decode recent messages: %w", err)
	}

	log.Printf("Retrieved %d recent messages for user: %s", len(records), userID)
	return records, nil
}

// GetMessageCount returns the total number of messages for a user
func (s *SMSService) GetMessageCount(ctx context.Context, userID string) (int64, error) {
	collection := db.GetCollection()

	queryCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	filter := bson.M{"user_id": userID}
	count, err := collection.CountDocuments(queryCtx, filter)
	if err != nil {
		return 0, fmt.Errorf("failed to count messages: %w", err)
	}

	return count, nil
}
