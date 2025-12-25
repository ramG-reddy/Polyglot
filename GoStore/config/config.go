package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

// Config holds all configuration for the SMS Store service
type Config struct {
	// Server Configuration
	ServerPort string

	// MongoDB Configuration
	MongoURI      string
	MongoDatabase string
	MongoUser     string
	MongoPassword string

	// Kafka Configuration
	KafkaBrokers []string
	KafkaTopic   string
	KafkaGroupID string
}

var AppConfig *Config

// Load reads configuration from environment variables with sensible defaults
func Load() (*Config, error) {
	log.Println("Loading configuration from environment variables...")

	config := &Config{
		ServerPort:    getEnv("GO_SERVICE_PORT", "8090"),
		MongoDatabase: getEnv("MONGO_DATABASE", "sms_store"),
		MongoUser:     getEnv("MONGO_APP_USER", "smsapp"),
		MongoPassword: getEnv("MONGO_APP_PASSWORD", "smsapp123"),
		KafkaTopic:    getEnv("KAFKA_TOPIC", "sms.events"),
		KafkaGroupID:  getEnv("KAFKA_GROUP_ID", "sms-store-consumer-group"),
	}

	// Build MongoDB connection URI
	mongoHost := getEnv("MONGO_HOST", "mongodb")
	mongoPort := getEnv("MONGO_PORT", "27017")
	config.MongoURI = fmt.Sprintf("mongodb://%s:%s@%s:%s/%s?authSource=%s",
		config.MongoUser,
		config.MongoPassword,
		mongoHost,
		mongoPort,
		config.MongoDatabase,
		config.MongoDatabase,
	)

	// Parse Kafka brokers (comma-separated list)
	kafkaBrokerList := getEnv("KAFKA_BROKERS", "kafka:9092")
	config.KafkaBrokers = []string{kafkaBrokerList}

	// Validate required configuration
	if err := config.validate(); err != nil {
		return nil, fmt.Errorf("configuration validation failed: %w", err)
	}

	AppConfig = config
	log.Printf("Configuration loaded successfully: Server Port=%s, Kafka Topic=%s, MongoDB=%s",
		config.ServerPort, config.KafkaTopic, config.MongoDatabase)

	return config, nil
}

// validate checks that all required configuration values are present
func (c *Config) validate() error {
	if c.ServerPort == "" {
		return fmt.Errorf("server port is required")
	}
	if c.MongoURI == "" {
		return fmt.Errorf("MongoDB URI is required")
	}
	if c.MongoDatabase == "" {
		return fmt.Errorf("MongoDB database name is required")
	}
	if len(c.KafkaBrokers) == 0 {
		return fmt.Errorf("at least one Kafka broker is required")
	}
	if c.KafkaTopic == "" {
		return fmt.Errorf("Kafka topic is required")
	}
	if c.KafkaGroupID == "" {
		return fmt.Errorf("Kafka group ID is required")
	}
	return nil
}

// getEnv retrieves an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsInt retrieves an environment variable as integer or returns default
func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		log.Printf("Warning: Invalid integer value for %s: %s, using default: %d", key, valueStr, defaultValue)
		return defaultValue
	}
	return value
}
