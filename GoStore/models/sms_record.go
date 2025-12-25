package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// SMSRecord represents a stored SMS message record in MongoDB
type SMSRecord struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID      string             `bson:"user_id" json:"user_id"`
	PhoneNumber string             `bson:"phone_number" json:"phone_number"`
	Message     string             `bson:"message" json:"message"`
	Status      string             `bson:"status" json:"status"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
}

// KafkaEvent represents the event consumed from Kafka topic
// This matches the Java KafkaEvent structure but uses Go types
type KafkaEvent struct {
	EventID     string `json:"eventId"`
	UserID      string `json:"userId"`
	PhoneNumber string `json:"phoneNumber"`
	Message     string `json:"message"`
	Status      string `json:"status"`
	CreatedAt   string `json:"createdAt"` // ISO-8601 format from Java (no timezone)
}

// ToSMSRecord converts a KafkaEvent to an SMSRecord for MongoDB storage
// Handles timestamp conversion from Java ISO-8601 (no TZ) to Go time.Time (UTC)
func (k *KafkaEvent) ToSMSRecord() (*SMSRecord, error) {
	// Parse Java LocalDateTime format (ISO-8601 without timezone)
	// Java sends: "2025-12-25T10:30:00"
	// We need to parse it and treat it as UTC
	createdAt, err := parseJavaLocalDateTime(k.CreatedAt)
	if err != nil {
		// If parsing fails, use current time as fallback
		createdAt = time.Now().UTC()
	}

	return &SMSRecord{
		UserID:      k.UserID,
		PhoneNumber: k.PhoneNumber,
		Message:     k.Message,
		Status:      k.Status,
		CreatedAt:   createdAt,
	}, nil
}

// parseJavaLocalDateTime parses Java LocalDateTime (ISO-8601 without timezone)
// and returns a Go time.Time in UTC
func parseJavaLocalDateTime(timestamp string) (time.Time, error) {
	// Java LocalDateTime format: "2025-12-25T10:30:00"
	// We need to append 'Z' to parse it as UTC
	layouts := []string{
		"2006-01-02T15:04:05",        // Without seconds fraction
		"2006-01-02T15:04:05.999999", // With microseconds
		time.RFC3339,                 // Fallback if already has timezone
	}

	var lastErr error
	for _, layout := range layouts {
		t, err := time.Parse(layout, timestamp)
		if err == nil {
			// Successfully parsed, return as UTC
			return t.UTC(), nil
		}
		lastErr = err
	}

	// If all parsing attempts failed, return the last error
	return time.Time{}, lastErr
}

// ListMessagesResponse is not needed as we return []SMSRecord directly
// The JSON marshaling will handle the array format automatically
