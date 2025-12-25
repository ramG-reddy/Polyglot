# Polyglot SMS Service - Test Demonstration Guide

This document provides step-by-step test scripts to demonstrate all features of the Polyglot SMS Service.

## Table of Contents

1. [Prerequisites](#prerequisites)
2. [Test 1: Verify All Services Are Running](#test-1-verify-all-services-are-running)
3. [Test 2: Send SMS Messages (Java Service)](#test-2-send-sms-messages-java-service)
   - [Test 2.1: Send Valid SMS](#test-21-send-valid-sms)
   - [Test 2.2: Send Multiple Messages](#test-22-send-multiple-messages)
   - [Test 2.3: Test Invalid Phone Number](#test-23-test-invalid-phone-number)
   - [Test 2.4: Test Empty Message](#test-24-test-empty-message)
   - [Test 2.5: Test Very Long Message](#test-25-test-very-long-message)
4. [Test 3: Blocked Users Feature](#test-3-blocked-users-feature)
   - [Test 3.1: Check Current Blocked Users](#test-31-check-current-blocked-users)
   - [Test 3.2: Send SMS to Blocked User](#test-32-send-sms-to-blocked-user)
   - [Test 3.3: Add Custom Blocked User](#test-33-add-custom-blocked-user)
   - [Test 3.4: Verify Block Works](#test-34-verify-block-works)
5. [Test 4: Retrieve Messages (Go Service)](#test-4-retrieve-messages-go-service)
   - [Test 4.1: Retrieve Messages for Specific User](#test-41-retrieve-messages-for-specific-user)
   - [Test 4.2: Retrieve Messages for Multiple Users](#test-42-retrieve-messages-for-multiple-users)
   - [Test 4.3: Retrieve Messages for Non-Existent User](#test-43-retrieve-messages-for-non-existent-user)
   - [Test 4.4: Test Invalid User ID Format](#test-44-test-invalid-user-id-format)
6. [Test 5: Verify Data in MongoDB](#test-5-verify-data-in-mongodb)
   - [Test 5.1: Connect to MongoDB and View All Records](#test-51-connect-to-mongodb-and-view-all-records)
   - [Test 5.2: Verify Record Structure](#test-52-verify-record-structure)
7. [Test 6: Verify Kafka Message Flow](#test-6-verify-kafka-message-flow)
   - [Test 6.1: Check Kafka Topics](#test-61-check-kafka-topics)
   - [Test 6.2: View Kafka Topic Details](#test-62-view-kafka-topic-details)
   - [Test 6.3: Consume Messages from Kafka Topic](#test-63-consume-messages-from-kafka-topic-real-time-monitoring)
   - [Test 6.4: Send SMS and Monitor Kafka in Real-time](#test-64-send-sms-and-monitor-kafka-in-real-time)
8. [Test 7: Verify Redis Block List](#test-7-verify-redis-block-list)
   - [Test 7.1: View All Blocked Users](#test-71-view-all-blocked-users)
   - [Test 7.2: Check if User is Blocked](#test-72-check-if-user-is-blocked)
   - [Test 7.3: Add New Blocked User](#test-73-add-new-blocked-user)
   - [Test 7.4: Remove User from Block List](#test-74-remove-user-from-block-list)
   - [Test 7.5: View All Redis Keys](#test-75-view-all-redis-keys)
9. [Test 8: End-to-End Workflow Test](#test-8-end-to-end-workflow-test)
10. [Test 9: Load Testing](#test-9-load-testing)
    - [Test 9.1: Send 50 Messages Rapidly](#test-91-send-50-messages-rapidly)
    - [Test 9.2: Verify All Messages Were Stored](#test-92-verify-all-messages-were-stored)
11. [Test 10: Error Scenarios](#test-10-error-scenarios)
    - [Test 10.1: Test with Malformed JSON](#test-101-test-with-malformed-json)
    - [Test 10.2: Test Missing Required Fields](#test-102-test-missing-required-fields)
    - [Test 10.3: Test Wrong Content-Type](#test-103-test-wrong-content-type)
12. [Test 11: Service Health Checks](#test-11-service-health-checks)
    - [Test 11.1: Check Java Service Health](#test-111-check-java-service-health)
    - [Test 11.2: Check Go Service Health](#test-112-check-go-service-health)
    - [Test 11.3: Check All Docker Health Status](#test-113-check-all-docker-health-status)
13. [Test 12: Clean Up and Reset](#test-12-clean-up-and-reset)
    - [Test 12.1: Clear All SMS Records](#test-121-clear-all-sms-records)
    - [Test 12.2: Reset Redis Block List](#test-122-reset-redis-block-list)
    - [Test 12.3: Restart All Services](#test-123-restart-all-services)
    - [Test 12.4: Stop All Services](#test-124-stop-all-services)
    - [Test 12.5: Complete Cleanup (Remove Volumes)](#test-125-complete-cleanup-remove-volumes)
14. [Test 13: Monitor Logs in Real-Time](#test-13-monitor-logs-in-real-time)
    - [Test 13.1: Follow All Service Logs](#test-131-follow-all-service-logs)
    - [Test 13.2: Follow Specific Service](#test-132-follow-specific-service)
15. [Quick Test Script](#quick-test-script)
16. [Expected Results Summary](#expected-results-summary)
17. [Troubleshooting](#troubleshooting)
18. [Test Completion Checklist](#test-completion-checklist)

---

## Prerequisites

Ensure Docker Desktop is running and all services are up:
```powershell
docker compose up -d
```

Wait for all services to be healthy (takes ~30-60 seconds):
```powershell
docker compose ps
```

---

## Test 1: Verify All Services Are Running

### Check Service Health Status
```powershell
docker compose ps
```

**Expected Output:** All services should show "Up (healthy)" status.

---

## Test 2: Send SMS Messages (Java Service)

### Test 2.1: Send Valid SMS
```powershell
curl -X POST http://localhost:8080/v0/sms/send `
  -H "Content-Type: application/json" `
  -d '{\"phoneNumber\": \"+1234567890\", \"message\": \"Hello from Polyglot SMS Service!\"}'
```

**Expected Output:**
```json
{
  "status": "SUCCESS",
  "message": "SMS sent successfully",
  "timestamp": "2025-12-25T..."
}
```
or
```json
{
  "status": "FAILED",
  "message": "3rd party API failed",
  "timestamp": "2025-12-25T..."
}
```

### Test 2.2: Send Multiple Messages
```powershell
# Send 5 messages to different users
for ($i=1; $i -le 5; $i++) {
    Write-Host "`n--- Sending SMS $i ---"
    curl -X POST http://localhost:8080/v0/sms/send `
      -H "Content-Type: application/json" `
      -d "{\"phoneNumber\": \"+123456789$i\", \"message\": \"Test message number $i\"}"
    Start-Sleep -Seconds 1
}
```

### Test 2.3: Test Invalid Phone Number
```powershell
curl -X POST http://localhost:8080/v0/sms/send `
  -H "Content-Type: application/json" `
  -d '{\"phoneNumber\": \"invalid\", \"message\": \"This should fail validation\"}'
```

**Expected Output:** Validation error (400 Bad Request)

### Test 2.4: Test Empty Message
```powershell
curl -X POST http://localhost:8080/v0/sms/send `
  -H "Content-Type: application/json" `
  -d '{\"phoneNumber\": \"+1234567890\", \"message\": \"\"}'
```

**Expected Output:** Validation error (400 Bad Request)

### Test 2.5: Test Very Long Message
```powershell
$longMessage = "A" * 500
curl -X POST http://localhost:8080/v0/sms/send `
  -H "Content-Type: application/json" `
  -d "{\"phoneNumber\": \"+1234567890\", \"message\": \"$longMessage\"}"
```

**Expected Output:** Validation error (message too long)

---

## Test 3: Blocked Users Feature

### Test 3.1: Check Current Blocked Users
```powershell
docker exec -it polyglot-redis redis-cli SMEMBERS blocked_users
```

**Expected Output:** List of blocked phone numbers (e.g., +1111111111, +2222222222, +3333333333)

### Test 3.2: Send SMS to Blocked User
```powershell
curl -X POST http://localhost:8080/v0/sms/send `
  -H "Content-Type: application/json" `
  -d '{\"phoneNumber\": \"+1111111111\", \"message\": \"This should be blocked\"}'
```

**Expected Output:**
```json
{
  "status": "BLOCKED",
  "message": "User is blocked",
  "timestamp": "2025-12-25T..."
}
```

### Test 3.3: Add Custom Blocked User
```powershell
docker exec -it polyglot-redis redis-cli SADD blocked_users "+9999999999"
```

### Test 3.4: Verify Block Works
```powershell
curl -X POST http://localhost:8080/v0/sms/send `
  -H "Content-Type: application/json" `
  -d '{\"phoneNumber\": \"+9999999999\", \"message\": \"Should be blocked\"}'
```

**Expected Output:** Status "blocked"

---

## Test 4: Retrieve Messages (Go Service)

### Test 4.1: Retrieve Messages for Specific User
```powershell
curl http://localhost:8090/v0/user/+1234567890/messages
```

**Expected Output:** JSON array of SMS records for that user

### Test 4.2: Retrieve Messages for Multiple Users
```powershell
# Check messages for first 5 test users
for ($i=1; $i -le 5; $i++) {
    Write-Host "`n--- Messages for +123456789$i ---"
    curl http://localhost:8090/v0/user/+123456789$i/messages
    Start-Sleep -Seconds 1
}
```

### Test 4.3: Retrieve Messages for Non-Existent User
```powershell
curl http://localhost:8090/v0/user/+0000000000/messages
```

**Expected Output:** Empty array `[]`

### Test 4.4: Test Invalid User ID Format
```powershell
curl http://localhost:8090/v0/user/invalid/messages
```

**Expected Output:** Empty array (no validation error, just no results)

---

## Test 5: Verify Data in MongoDB

### Test 5.1: Connect to MongoDB and View All Records
```powershell
docker exec -it polyglot-mongodb mongosh -u smsapp -p smsapp123 --authenticationDatabase sms_store
```

Inside MongoDB shell:
```javascript
// Switch to database
use sms_store

// Count total records
db.sms_records.countDocuments()

// View all records (limit 10)
db.sms_records.find().limit(10).pretty()

// View records for specific user
db.sms_records.find({user_id: "+1234567890"}).pretty()

// View only successful messages
db.sms_records.find({status: "success"}).pretty()

// View only failed messages
db.sms_records.find({status: "failed"}).pretty()

// View only blocked messages
db.sms_records.find({status: "blocked"}).pretty()

// Get statistics by status
db.sms_records.aggregate([
  {$group: {_id: "$status", count: {$sum: 1}}}
])

// View recent messages (sorted by timestamp)
db.sms_records.find().sort({created_at: -1}).limit(10).pretty()

// Check indexes
db.sms_records.getIndexes()

// Exit
exit
```

### Test 5.2: Verify Record Structure
```powershell
docker exec -it polyglot-mongodb mongosh -u smsapp -p smsapp123 --authenticationDatabase sms_store --eval "db.getSiblingDB('sms_store').sms_records.findOne()"
```

**Expected Output:** Document with fields: `_id`, `user_id`, `message`, `status`, `created_at`

---

## Test 6: Verify Kafka Message Flow

### Test 6.1: Check Kafka Topics
```powershell
docker exec -it polyglot-kafka kafka-topics --bootstrap-server localhost:9092 --list
```

**Expected Output:** Should include `sms.events` topic

### Test 6.2: View Kafka Topic Details
```powershell
docker exec -it polyglot-kafka kafka-topics --bootstrap-server localhost:9092 --describe --topic sms.events
```

### Test 6.3: Consume Messages from Kafka Topic (Real-time Monitoring)
```powershell
docker exec -it polyglot-kafka kafka-console-consumer --bootstrap-server localhost:9092 --topic sms.events --from-beginning
```

**Note:** Press `Ctrl+C` to stop. You should see JSON messages.

### Test 6.4: Send SMS and Monitor Kafka in Real-time

**Terminal 1:** Start Kafka consumer
```powershell
docker exec -it polyglot-kafka kafka-console-consumer --bootstrap-server localhost:9092 --topic sms.events
```

**Terminal 2:** Send SMS
```powershell
curl -X POST http://localhost:8080/v0/sms/send `
  -H "Content-Type: application/json" `
  -d '{\"phoneNumber\": \"+1234567890\", \"message\": \"Real-time Kafka test\"}'
```

**Expected:** Terminal 1 should immediately show the Kafka event.

---

## Test 7: Verify Redis Block List

### Test 7.1: View All Blocked Users
```powershell
docker exec -it polyglot-redis redis-cli SMEMBERS blocked_users
```

### Test 7.2: Check if User is Blocked
```powershell
docker exec -it polyglot-redis redis-cli SISMEMBER blocked_users "+1111111111"
```

**Expected Output:** `1` (true) if blocked, `0` (false) if not

### Test 7.3: Add New Blocked User
```powershell
docker exec -it polyglot-redis redis-cli SADD blocked_users "+5555555555"
```

### Test 7.4: Remove User from Block List
```powershell
docker exec -it polyglot-redis redis-cli SREM blocked_users "+5555555555"
```

### Test 7.5: View All Redis Keys
```powershell
docker exec -it polyglot-redis redis-cli KEYS "*"
```

---

## Test 8: End-to-End Workflow Test

This test demonstrates the complete flow from sending an SMS to retrieving it.

```powershell
Write-Host "`n=== POLYGLOT SMS SERVICE - END-TO-END TEST ===`n"

# Step 1: Send SMS
Write-Host "Step 1: Sending SMS to +1234567890..."
$response = curl -X POST http://localhost:8080/v0/sms/send `
  -H "Content-Type: application/json" `
  -d '{\"phoneNumber\": \"+1234567890\", \"message\": \"End-to-end test message\"}'
Write-Host "Response: $response`n"
Start-Sleep -Seconds 2

# Step 2: Wait for Kafka processing
Write-Host "Step 2: Waiting for Kafka consumer to process message..."
Start-Sleep -Seconds 3

# Step 3: Retrieve messages from Go service
Write-Host "Step 3: Retrieving messages for +1234567890..."
$messages = curl http://localhost:8090/v0/user/+1234567890/messages
Write-Host "Messages: $messages`n"

# Step 4: Verify in MongoDB
Write-Host "Step 4: Verifying in MongoDB..."
docker exec -it polyglot-mongodb mongosh -u smsapp -p smsapp123 --authenticationDatabase sms_store --eval "db.getSiblingDB('sms_store').sms_records.find({user_id: '+1234567890'}).sort({created_at: -1}).limit(1).pretty()"

Write-Host "`n=== TEST COMPLETE ===`n"
```

---

## Test 9: Load Testing

### Test 9.1: Send 50 Messages Rapidly
```powershell
Write-Host "Sending 50 messages..."
for ($i=1; $i -le 50; $i++) {
    $user = "+10000000" + ([string]($i % 10)).PadLeft(2, '0')
    curl -X POST http://localhost:8080/v0/sms/send `
      -H "Content-Type: application/json" `
      -d "{\"phoneNumber\": \"$user\", \"message\": \"Load test message $i\"}" | Out-Null
    if ($i % 10 -eq 0) {
        Write-Host "Sent $i messages..."
    }
}
Write-Host "Load test complete!"
```

### Test 9.2: Verify All Messages Were Stored
```powershell
Start-Sleep -Seconds 5
docker exec -it polyglot-mongodb mongosh -u smsapp -p smsapp123 --authenticationDatabase sms_store --eval "db.getSiblingDB('sms_store').sms_records.countDocuments()"
```

---

## Test 10: Error Scenarios

### Test 10.1: Test with Malformed JSON
```powershell
curl -X POST http://localhost:8080/v0/sms/send `
  -H "Content-Type: application/json" `
  -d '{invalid json}'
```

**Expected Output:** 400 Bad Request

### Test 10.2: Test Missing Required Fields
```powershell
curl -X POST http://localhost:8080/v0/sms/send `
  -H "Content-Type: application/json" `
  -d '{\"phoneNumber\": \"+1234567890\"}'
```

**Expected Output:** 400 Bad Request (missing message)

### Test 10.3: Test Wrong Content-Type
```powershell
curl -X POST http://localhost:8080/v0/sms/send `
  -H "Content-Type: text/plain" `
  -d 'plain text data'
```

**Expected Output:** 415 Unsupported Media Type

---

## Test 11: Service Health Checks

### Test 11.1: Check Java Service Health
```powershell
curl http://localhost:8080/actuator/health
```

**Expected Output:** `{"status":"UP",...}`

### Test 11.2: Check Go Service Health
```powershell
curl http://localhost:8090/health
```

**Expected Output:** `{"status":"ok","database":"connected",...}`

### Test 11.3: Check All Docker Health Status
```powershell
docker ps --format "table {{.Names}}\t{{.Status}}"
```

---

## Test 12: Clean Up and Reset

### Test 12.1: Clear All SMS Records
```powershell
docker exec -it polyglot-mongodb mongosh -u smsapp -p smsapp123 --authenticationDatabase sms_store --eval "db.getSiblingDB('sms_store').sms_records.deleteMany({})"
```

### Test 12.2: Reset Redis Block List
```powershell
docker exec -it polyglot-redis redis-cli DEL blocked_users
docker exec -it polyglot-redis redis-cli SADD blocked_users "+1111111111" "+2222222222" "+3333333333"
```

### Test 12.3: Restart All Services
```powershell
docker compose restart
```

### Test 12.4: Stop All Services
```powershell
docker compose stop
```

### Test 12.5: Complete Cleanup (Remove Volumes)
```powershell
docker compose down -v
```

**Warning:** This will delete all data!

---

## Test 13: Monitor Logs in Real-Time

### Test 13.1: Follow All Service Logs
```powershell
docker compose logs -f
```

### Test 13.2: Follow Specific Service
```powershell
# Java service
docker compose logs -f sms-sender

# Go service
docker compose logs -f sms-store

# Kafka
docker compose logs -f kafka
```

---

## Quick Test Script

Save this as `quick-test.ps1` for rapid testing:

```powershell
# Quick Test Script for Polyglot SMS Service

Write-Host "`n=== QUICK TEST STARTED ===`n" -ForegroundColor Green

# Test 1: Send successful message
Write-Host "Test 1: Sending valid SMS..." -ForegroundColor Yellow
curl -X POST http://localhost:8080/v0/sms/send -H "Content-Type: application/json" -d '{\"phoneNumber\": \"+1234567890\", \"message\": \"Test message\"}'
Start-Sleep -Seconds 2

# Test 2: Send to blocked user
Write-Host "`nTest 2: Sending to blocked user..." -ForegroundColor Yellow
curl -X POST http://localhost:8080/v0/sms/send -H "Content-Type: application/json" -d '{\"phoneNumber\": \"+1111111111\", \"message\": \"Should be blocked\"}'
Start-Sleep -Seconds 2

# Test 3: Retrieve messages
Write-Host "`nTest 3: Retrieving messages..." -ForegroundColor Yellow
curl http://localhost:8090/v0/user/+1234567890/messages
Start-Sleep -Seconds 2

# Test 4: Check MongoDB
Write-Host "`nTest 4: MongoDB record count..." -ForegroundColor Yellow
docker exec polyglot-mongodb mongosh -u smsapp -p smsapp123 --authenticationDatabase sms_store --quiet --eval "db.getSiblingDB('sms_store').sms_records.countDocuments()"

# Test 5: Check Redis
Write-Host "`nTest 5: Redis blocked users..." -ForegroundColor Yellow
docker exec polyglot-redis redis-cli SMEMBERS blocked_users

Write-Host "`n=== QUICK TEST COMPLETED ===`n" -ForegroundColor Green
```

Run it with:
```powershell
.\quick-test.ps1
```

---

## Expected Results Summary

- ✅ Java service accepts SMS requests on port 8080
- ✅ Blocked users are rejected immediately via Redis
- ✅ Valid requests produce Kafka events
- ✅ Go service consumes Kafka events and stores in MongoDB
- ✅ Messages can be retrieved via Go service on port 8090
- ✅ All data persists across restarts
- ✅ Validation works correctly for all input types
- ✅ Health checks return positive status

---

## Troubleshooting

If any test fails:

1. **Check service logs:**
   ```powershell
   docker compose logs <service-name>
   ```

2. **Verify all services are healthy:**
   ```powershell
   docker compose ps
   ```

3. **Restart specific service:**
   ```powershell
   docker compose restart <service-name>
   ```

4. **Complete restart:**
   ```powershell
   docker compose down
   docker compose up -d
   ```

5. **Check network connectivity:**
   ```powershell
   docker network inspect polyglot-network
   ```

---

## Test Completion Checklist

- [ ] All services start successfully
- [ ] Java service sends SMS successfully
- [ ] Blocked users are rejected
- [ ] Kafka events are produced
- [ ] Go service consumes events
- [ ] Messages are stored in MongoDB
- [ ] Messages can be retrieved
- [ ] Validation works correctly
- [ ] Health checks pass
- [ ] Load testing works without errors
- [ ] Data persists across restarts

---

**End of Test Documentation**
