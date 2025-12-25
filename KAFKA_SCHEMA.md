# Kafka Topics and Message Schema

This document describes the Kafka topics, message schemas, and event contracts used in the Polyglot SMS Service.

## Overview

The system uses Apache Kafka as a message broker to enable asynchronous communication between the Java SMS Sender service (producer) and the Go SMS Store service (consumer).

**Architecture Flow:**
```
[Java SMS Sender] --> Kafka Topic: sms.events --> [Go SMS Store] --> MongoDB
```

---

## Topics

### `sms.events`

**Purpose**: Streams SMS transaction events from the sender service to the store service

**Configuration:**
- **Partitions**: 3 (configured in docker-compose.yml)
- **Replication Factor**: 1 (single broker setup)
- **Retention**: 168 hours (7 days)
- **Compression**: None (for simplicity with polyglot consumers)
- **Auto-creation**: Enabled

**Producer**: Java SMS Sender Service
**Consumer Group**: `sms-store-consumer-group` (Go SMS Store Service)

---

## Message Schema

### SMS Event Message

**Topic**: `sms.events`  
**Format**: JSON  
**Key**: `null` (messages are not keyed)  
**Value**: SMS Event JSON object

#### JSON Schema

```json
{
  "eventId": "string (UUID)",
  "userId": "string (E.164 phone number)",
  "phoneNumber": "string (E.164 phone number)",
  "message": "string (1-160 characters)",
  "status": "string (SUCCESS|FAILED|BLOCKED)",
  "createdAt": "string (ISO-8601 datetime)"
}
```

#### Field Descriptions

| Field | Type | Required | Description | Example |
|-------|------|----------|-------------|---------|
| `eventId` | String (UUID) | Yes | Unique identifier for this event | `"7a61ec00-3391-47ac-8420-38b3537f9a72"` |
| `userId` | String | Yes | User identifier (same as phoneNumber in this system) | `"+1234567890"` |
| `phoneNumber` | String | Yes | Destination phone number in E.164 format | `"+1234567890"` |
| `message` | String | Yes | SMS message content (1-160 chars) | `"Hello from Polyglot SMS!"` |
| `status` | String (Enum) | Yes | Status of SMS operation | `"SUCCESS"` |
| `createdAt` | String (ISO-8601) | Yes | Timestamp when event was created | `"2025-12-26T10:30:45"` |

#### Status Values

| Status | Description | When Used |
|--------|-------------|-----------|
| `SUCCESS` | SMS was successfully sent by vendor API | Vendor API returned 200 OK |
| `FAILED` | SMS sending failed at vendor API | Vendor API returned error (simulated random failures) |
| `BLOCKED` | User is on the block list, SMS was rejected | Phone number found in Redis block list |

#### Timestamp Format

The `createdAt` field uses ISO-8601 format without timezone information:
- **Format**: `yyyy-MM-dd'T'HH:mm:ss`
- **Example**: `2025-12-26T14:23:10`
- **Interpretation**: Treated as UTC time by the Go consumer

---

## Message Examples

### Successful SMS Event

```json
{
  "eventId": "a1b2c3d4-5678-90ab-cdef-1234567890ab",
  "userId": "+1234567890",
  "phoneNumber": "+1234567890",
  "message": "Your verification code is 123456",
  "status": "SUCCESS",
  "createdAt": "2025-12-26T10:30:45"
}
```

### Failed SMS Event

```json
{
  "eventId": "e5f6g7h8-9012-34ij-klmn-5678901234op",
  "userId": "+9876543210",
  "phoneNumber": "+9876543210",
  "message": "This SMS failed to send",
  "status": "FAILED",
  "createdAt": "2025-12-26T10:31:12"
}
```

### Blocked User Event

```json
{
  "eventId": "q9r8s7t6-u5v4-w3x2-y1z0-abcdef123456",
  "userId": "+1111111111",
  "phoneNumber": "+1111111111",
  "message": "This user is blocked",
  "status": "BLOCKED",
  "createdAt": "2025-12-26T10:32:05"
}
```

---

## Producer Implementation

### Java (Spring Kafka)

**Location**: `JavaSender/src/main/java/com/sms/sender/kafka/SmsKafkaProducer.java`

**Configuration**:
```java
// Key Serializer: StringSerializer
// Value Serializer: JsonSerializer
// Compression: None
// Acknowledgments: all (wait for all replicas)
// Idempotence: Enabled
```

**Production Flow**:
1. Convert `SmsRequest` to `KafkaEvent` model
2. Serialize to JSON using Jackson
3. Send synchronously to `sms.events` topic
4. Wait for acknowledgment from Kafka broker
5. Return success/failure to caller

**Message Key**: `null` (round-robin distribution across partitions)

---

## Consumer Implementation

### Go (segmentio/kafka-go)

**Location**: `GoStore/kafka/consumer.go`

**Configuration**:
```go
// Consumer Group: sms-store-consumer-group
// Start Offset: FirstOffset (reads from beginning)
// Commit Interval: 1 second
// Max Wait: 500ms
// Auto-commit: Disabled (manual commit after processing)
```

**Consumption Flow**:
1. Fetch message from Kafka topic
2. Deserialize JSON to `KafkaEvent` struct
3. Convert to `SMSRecord` model
4. Persist to MongoDB
5. Commit offset to Kafka
6. Log success/failure

**Error Handling**:
- Parse errors: Skip message and log error
- Database errors: Retry (message not committed)
- Timeout: Continue to next message

---

## Data Transformations

### Java → Kafka

**Java Model** (`KafkaEvent.java`):
```java
@JsonProperty("createdAt")
@JsonFormat(shape = JsonFormat.Shape.STRING, pattern = "yyyy-MM-dd'T'HH:mm:ss")
private LocalDateTime createdAt;
```

**JSON Output**:
```json
"createdAt": "2025-12-26T10:30:45"
```

### Kafka → Go

**Go Model** (`models/sms_record.go`):
```go
type KafkaEvent struct {
    CreatedAt string `json:"createdAt"` // ISO-8601 format
}
```

**Conversion**:
```go
// Parse ISO-8601 string to time.Time
createdAt, _ := time.Parse("2006-01-02T15:04:05", event.CreatedAt)
```

### Go → MongoDB

**Go Model** (`SMSRecord`):
```go
type SMSRecord struct {
    CreatedAt time.Time `bson:"created_at" json:"created_at"`
}
```

**BSON Storage**: Stored as MongoDB ISODate type

---

## Monitoring and Operations

### View Messages in Topic

```bash
# From beginning
docker exec -it polyglot-kafka kafka-console-consumer \
  --bootstrap-server localhost:9092 \
  --topic sms.events \
  --from-beginning

# Real-time monitoring
docker exec -it polyglot-kafka kafka-console-consumer \
  --bootstrap-server localhost:9092 \
  --topic sms.events
```

### Topic Details

```bash
# List all topics
docker exec -it polyglot-kafka kafka-topics \
  --bootstrap-server localhost:9092 \
  --list

# Describe topic
docker exec -it polyglot-kafka kafka-topics \
  --bootstrap-server localhost:9092 \
  --describe \
  --topic sms.events
```

### Consumer Group Status

```bash
# List consumer groups
docker exec -it polyglot-kafka kafka-consumer-groups \
  --bootstrap-server localhost:9092 \
  --list

# Describe consumer group
docker exec -it polyglot-kafka kafka-consumer-groups \
  --bootstrap-server localhost:9092 \
  --describe \
  --group sms-store-consumer-group
```

### Check Lag

```bash
# View consumer lag
docker exec -it polyglot-kafka kafka-consumer-groups \
  --bootstrap-server localhost:9092 \
  --describe \
  --group sms-store-consumer-group \
  --members
```

---

## Schema Evolution

### Current Version: v1.0

**Breaking Changes** require:
1. Update both producer and consumer
2. Create new topic or handle both formats
3. Coordinate deployment

**Non-Breaking Changes** (backward compatible):
- Adding optional fields to the end
- Adding new status values (if consumer handles unknown values gracefully)

### Future Considerations

1. **Schema Registry**: Consider Confluent Schema Registry or similar for schema management
2. **Versioning**: Add `schemaVersion` field to messages
3. **Avro/Protobuf**: Consider binary formats for better performance
4. **Dead Letter Queue**: Add DLQ topic for failed message processing

---

## Performance Characteristics

### Throughput

- **Producer**: Synchronous sends (blocking) - ~100-500 msg/sec
- **Consumer**: Asynchronous batch processing - ~1000+ msg/sec
- **Bottleneck**: MongoDB write operations

### Latency

- **Kafka Produce**: ~5-20ms (local network)
- **Kafka Consume**: ~10-50ms (including MongoDB write)
- **End-to-End**: ~50-200ms (HTTP request → MongoDB storage)

### Scalability

- **Horizontal**: Can add more consumer instances (up to 3 with current partitions)
- **Vertical**: Increase Kafka broker resources
- **Partitioning**: Currently 3 partitions, can be increased

---

## Best Practices

### Producer Side (Java)

✅ **Do**:
- Use idempotent producer (enabled)
- Wait for all replica acknowledgments
- Log all produce attempts
- Handle produce failures gracefully
- Use JSON serialization for polyglot compatibility

❌ **Don't**:
- Send sensitive data unencrypted
- Produce without error handling
- Use blocking calls in async contexts
- Ignore failed produce operations

### Consumer Side (Go)

✅ **Do**:
- Use consumer groups for coordination
- Commit offsets only after successful processing
- Handle deserialization errors gracefully
- Log all consume attempts
- Implement retry logic for transient failures

❌ **Don't**:
- Auto-commit without processing
- Block consumption for long operations
- Ignore malformed messages silently
- Process same message multiple times (ensure idempotency)

---

## Troubleshooting

### Messages Not Being Consumed

1. Check consumer is running: `docker-compose logs sms-store`
2. Verify consumer group: `kafka-consumer-groups --describe`
3. Check topic exists: `kafka-topics --list`
4. Verify network connectivity between services

### Duplicate Messages

Possible causes:
- Consumer crashed before committing offset
- Network partition during commit
- Consumer rebalancing

Solution: Implement idempotency in consumer (check if record already exists)

### Message Format Errors

```
Error: failed to unmarshal Kafka event
```

Causes:
- Schema mismatch between producer and consumer
- Timestamp format incompatibility
- Incorrect JSON serialization

Solution: Verify `KafkaEvent` models match in both services

---

**Last Updated**: December 26, 2025  
**Schema Version**: 1.0
