# Polyglot Distributed SMS Service: Project Constitution

## 1. Project Goal
To refactor a monolithic notification system into a distributed, polyglot microservices architecture. The system comprises a Java-based **SMS Sender** (Gateway) and a GoLang-based **SMS Store** (Persistence), communicating via both synchronous HTTP and asynchronous Kafka events.

---

## 2. Architecture & Roles

### Service A: SMS Sender (The Gateway)
* **Language:** Java (Spring Boot).
* **Responsibility:**
    * Act as the public entry point for SMS requests.
    * Validate users against a **Block List** stored in Redis.
    * Interface with a (mocked) 3rd Party SMS Vendor.
    * Propagate metadata to Service B via Kafka (for logging) and/or HTTP (direct communication).

### Service B: SMS Store (The Vault)
* **Language:** GoLang (Standard Lib + Mongo Driver).
* **Responsibility:**
    * Ingest SMS records.
    * Persist data into **MongoDB** (Note: Architecture diagram mentions MySQL, but requirements strictly mandate MongoDB).
    * Serve retrieval APIs for user message history.

### Infrastructure
* **Message Broker:** Apache Kafka.
* **Cache:** Redis (for user block list).
* **Database:** MongoDB.
* **Containerization:** All services (Java SMS Sender, Go SMS Store, Kafka, Redis, MongoDB, Zookeeper) will run in Docker containers via docker-compose.

---

## 3. The Laws of Data Flow
The system must adhere to the following strict operational flow:
1.  **Ingestion:** Client POSTs request to Java Service.
2.  **Validation:** Java Service checks Redis. If blocked -> Stop.
3.  **Execution:** Java Service calls 3rd Party API (Mock: Randomly return SUCCESS/FAIL).
4.  **Logging (Async):** Java Service produces an event to Kafka.
5.  **Persistence:** Go Service consumes the Kafka event AND/OR accepts synchronous HTTP calls to store the record in MongoDB.
6.  **Retrieval:** Client GETs history from Go Service.

---

## 4. Coding Standards & Compliance
* **No Frameworks for Go HTTP:** The Go service must use the standard `net/http` library for routing, not Gin or Echo.
* **Polyglot Persistence:** Java owns Redis; Go owns MongoDB.
* **Error Handling:** Both services must handle timeouts and 3rd party failures gracefully.
* **Testing:** Core business logic must be unit tested in both languages.

---

## 5. Deliverables 
1. API Documentation: A simple README.md file detailing the endpoints for both 
services and instructions on how to run them locally. 
2. Demonstration: A script or documented steps demonstrating the full, end-to-end flow: 
Call the Java service, check the GoLang service's logs for the internal call, and finally 
retrieve the record using the GoLang service's history API.

---

## 6. Recommended Practices 
1. Code Structure: Organize code logically (handlers, services, models/structs). 
2. Error Handling: Implement robust logging and error handling, especially for inter-service communication timeouts or failures. 
3. Testing: Include basic Unit Tests for core business logic in both services. 
4. Best Practices: Follow language-specific best practices (Spring Boot conventions for Java and idiomatic Go for the Go service). 

---

## 7. Implementation Details & Decisions

### 7.1 Communication Architecture
**Decision**: Asynchronous Kafka-only between services (no direct HTTP between Java and Go services)
- **Rationale**: Decouples services, provides better fault tolerance, enables horizontal scaling
- **Implementation**: Java produces to `sms.events` topic; Go consumes from same topic
- **Trade-off**: Slightly higher latency (~50-200ms) but better reliability and scalability

### 7.2 Kafka Configuration
**Decision**: No compression, JSON serialization
- **Rationale**: Polyglot compatibility; Go consumer doesn't need Snappy codec
- **Implementation**: `compression.type=none` in producer config
- **Alternative Considered**: Snappy compression (rejected due to cross-language complexity)

### 7.3 Timestamp Handling
**Decision**: Java `LocalDateTime` serialized as ISO-8601 string
- **Rationale**: Cross-language compatibility; Go can parse string format easily
- **Implementation**: `@JsonFormat(pattern="yyyy-MM-dd'T'HH:mm:ss")` annotation
- **Trade-off**: String parsing overhead vs. array format incompatibility

### 7.4 Consumer Offset Management
**Decision**: Manual commit after successful MongoDB persistence
- **Rationale**: Ensures at-least-once delivery; prevents data loss
- **Implementation**: `StartOffset: kafka.FirstOffset` with manual commit
- **Trade-off**: Possible duplicate processing (requires idempotency)

### 7.5 MongoDB Indexing
**Decision**: Compound index on `(user_id, created_at)` plus individual indexes
- **Rationale**: Optimizes both user lookups and time-based queries
- **Implementation**: Automatic index creation on startup with existing index cleanup
- **Performance**: ~10ms query time for user message retrieval

### 7.6 Redis Block List
**Decision**: Redis Set data structure for blocked users
- **Rationale**: O(1) membership check; simple and efficient
- **Implementation**: `SISMEMBER blocked_users <phone>`
- **Initialization**: Populated with dummy data on Java service startup

### 7.7 Mock Vendor API
**Decision**: Configurable delay and failure rate
- **Rationale**: Realistic testing of timeout and error scenarios
- **Default Config**: 100-500ms delay, 30% failure rate
- **Environment Variables**: `APP_MOCK_VENDOR_MIN_DELAY_MS`, `APP_MOCK_VENDOR_FAILURE_RATE`

### 7.8 Health Checks
**Decision**: Docker healthchecks for all services with proper start periods
- **Rationale**: Ensures services are fully ready before marking as healthy
- **Implementation**: 
  - Java: `/actuator/health` endpoint
  - Go: `/health` endpoint
  - Infrastructure: Native commands (redis-cli, mongosh, kafka-broker-api-versions)

### 7.9 Validation Strategy
**Decision**: Jakarta Bean Validation in Java, manual validation in Go
- **Rationale**: Leverages Spring Boot's built-in validation framework
- **Java Validators**: 
  - Phone: `@Pattern(regexp="^\\+[1-9]\\d{1,14}$")`
  - Message: `@NotBlank`, `@Size(min=1, max=160)`
- **Go Validation**: Implicit through MongoDB queries (no explicit validation layer)

### 7.10 Dockerization Strategy
**Decision**: Multi-stage builds for both services
- **Rationale**: Smaller runtime images; separation of build and runtime dependencies
- **Java**: Maven build stage + JRE runtime (alpine)
- **Go**: Go build stage + Alpine runtime
- **Image Sizes**: Java ~400MB, Go ~20MB

---

## 8. Known Limitations & Future Improvements

### 8.1 Current Limitations
1. **No Message Deduplication**: Duplicate events possible on consumer restart
   - **Mitigation**: Can be addressed by checking MongoDB for existing eventId
2. **Single Kafka Broker**: No fault tolerance for broker failures
   - **Production**: Requires 3+ broker cluster with replication
3. **No TLS/Authentication**: Services communicate over plain text
   - **Production**: Requires SSL/TLS for all inter-service communication
4. **Mock Vendor**: Not a real SMS gateway
   - **Production**: Replace with actual vendor API (Twilio, AWS SNS, etc.)

### 8.2 Future Enhancements
1. **Schema Registry**: Add Confluent Schema Registry for Avro schemas
2. **Distributed Tracing**: Add OpenTelemetry for request tracing
3. **Metrics**: Add Prometheus metrics for monitoring
4. **Circuit Breaker**: Add circuit breaker pattern for vendor API calls
5. **Rate Limiting**: Add rate limiting to prevent API abuse
6. **Authentication**: Add JWT/OAuth2 authentication for APIs
7. **Idempotency Keys**: Add idempotency key support for duplicate prevention

---

## 9. Compliance & Verification

### 9.1 Constitution Compliance Checklist
- [x] Java service uses Spring Boot framework
- [x] Go service uses standard `net/http` library (no Gin/Echo)
- [x] Java owns Redis (block list management)
- [x] Go owns MongoDB (persistence layer)
- [x] Kafka for asynchronous event streaming
- [x] All services containerized via Docker Compose
- [x] Graceful error handling and timeout management
- [x] API documentation provided (README.md, test.md)
- [x] End-to-end demonstration scripts (test.md)

### 9.2 Data Flow Verification
1. ✅ **Ingestion**: Client → Java POST `/v0/sms/send`
2. ✅ **Validation**: Java checks Redis block list
3. ✅ **Execution**: Java calls mock vendor API (random success/failure)
4. ✅ **Logging**: Java produces Kafka event to `sms.events` topic
5. ✅ **Persistence**: Go consumes Kafka event → stores in MongoDB
6. ✅ **Retrieval**: Client → Go GET `/v0/user/{user_id}/messages`

---

**Last Updated**: December 26, 2025  
**Implementation Status**: Phase 4 Complete, Phase 5-6 In Progress