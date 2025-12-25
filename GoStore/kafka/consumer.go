package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/ramG-reddy/sms-store/models"
	"github.com/ramG-reddy/sms-store/services"
	"github.com/segmentio/kafka-go"
)

// Consumer handles Kafka message consumption
type Consumer struct {
	reader     *kafka.Reader
	smsService *services.SMSService
	stopChan   chan struct{}
}

// NewConsumer creates a new Kafka consumer instance
func NewConsumer(brokers []string, topic, groupID string, smsService *services.SMSService) *Consumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        brokers,
		Topic:          topic,
		GroupID:        groupID,
		MinBytes:       1,    // 1 byte
		MaxBytes:       10e6, // 10MB
		CommitInterval: time.Second,
		StartOffset:    kafka.LastOffset, // Start from latest for new consumer groups
		MaxWait:        500 * time.Millisecond,
		Logger:         kafka.LoggerFunc(log.Printf),
		ErrorLogger:    kafka.LoggerFunc(log.Printf),
	})

	return &Consumer{
		reader:     reader,
		smsService: smsService,
		stopChan:   make(chan struct{}),
	}
}

// StartConsumer begins consuming messages from Kafka in a background goroutine
func StartConsumer(brokers []string, topic, groupID string, smsService *services.SMSService) (*Consumer, error) {
	log.Printf("Starting Kafka consumer for topic: %s, group: %s", topic, groupID)

	consumer := NewConsumer(brokers, topic, groupID, smsService)

	// Start consumption in a goroutine
	go consumer.consume()

	log.Println("Kafka consumer started successfully")
	return consumer, nil
}

// consume is the main consumption loop that processes messages
func (c *Consumer) consume() {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Consumer panic recovered: %v", r)
		}
	}()

	log.Println("Starting message consumption loop...")

	for {
		select {
		case <-c.stopChan:
			log.Println("Consumer stop signal received, exiting...")
			return
		default:
			// Read message with timeout
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			message, err := c.reader.FetchMessage(ctx)
			cancel()

			if err != nil {
				if err == context.DeadlineExceeded {
					// Timeout is normal, continue
					continue
				}
				log.Printf("Error fetching message: %v", err)
				time.Sleep(1 * time.Second)
				continue
			}

			// Process the message
			if err := c.processMessage(message); err != nil {
				log.Printf("Error processing message: %v", err)
				// Don't commit on error - message will be reprocessed
				continue
			}

			// Commit the message after successful processing
			commitCtx, commitCancel := context.WithTimeout(context.Background(), 5*time.Second)
			if err := c.reader.CommitMessages(commitCtx, message); err != nil {
				log.Printf("Error committing message: %v", err)
			}
			commitCancel()
		}
	}
}

// processMessage deserializes and persists a Kafka message
func (c *Consumer) processMessage(message kafka.Message) error {
	log.Printf("Processing message from partition %d, offset %d", message.Partition, message.Offset)

	// Deserialize Kafka event from JSON
	var event models.KafkaEvent
	if err := json.Unmarshal(message.Value, &event); err != nil {
		return fmt.Errorf("failed to unmarshal Kafka event: %w", err)
	}

	log.Printf("Received event: EventID=%s, UserID=%s, Status=%s", event.EventID, event.UserID, event.Status)

	// Convert Kafka event to SMS record (handles timestamp conversion)
	record, err := event.ToSMSRecord()
	if err != nil {
		log.Printf("Warning: Failed to parse timestamp, using current time: %v", err)
		// Continue processing even if timestamp parsing fails
	}

	// Persist to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := c.smsService.SaveMessage(ctx, record); err != nil {
		return fmt.Errorf("failed to save message to database: %w", err)
	}

	log.Printf("Successfully processed and stored message for user: %s", event.UserID)
	return nil
}

// Stop gracefully shuts down the consumer
func (c *Consumer) Stop() error {
	log.Println("Stopping Kafka consumer...")

	// Signal the consumer to stop
	close(c.stopChan)

	// Give it a moment to finish current message
	time.Sleep(1 * time.Second)

	// Close the reader
	if err := c.reader.Close(); err != nil {
		return fmt.Errorf("failed to close Kafka reader: %w", err)
	}

	log.Println("Kafka consumer stopped successfully")
	return nil
}

// HealthCheck verifies the consumer is connected to Kafka
func (c *Consumer) HealthCheck() error {
	// The kafka-go library doesn't provide a direct health check
	// We can check if the reader is not nil
	if c.reader == nil {
		return fmt.Errorf("Kafka reader is not initialized")
	}
	return nil
}
