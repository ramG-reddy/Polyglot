# Environment Variables Documentation

This document describes all environment variables used in the Polyglot SMS Service.

## Java SMS Sender Service

### Spring Boot Configuration

| Variable Name | Default Value | Description | Required |
|--------------|---------------|-------------|----------|
| `SERVER_PORT` | `8080` | HTTP server port for the REST API | No |
| `SPRING_PROFILES_ACTIVE` | `docker` | Active Spring profile (local, docker, prod) | No |
| `LOGGING_LEVEL_COM_SMS_SENDER` | `INFO` | Logging level for the application (DEBUG, INFO, WARN, ERROR) | No |

### Kafka Configuration

| Variable Name | Default Value | Description | Required |
|--------------|---------------|-------------|----------|
| `SPRING_KAFKA_BOOTSTRAP_SERVERS` | `kafka:9092` | Comma-separated list of Kafka broker addresses | Yes |
| `SPRING_KAFKA_PRODUCER_KEY_SERIALIZER` | `StringSerializer` | Serializer class for Kafka message keys | No |
| `SPRING_KAFKA_PRODUCER_VALUE_SERIALIZER` | `JsonSerializer` | Serializer class for Kafka message values | No |
| `KAFKA_TOPIC` | `sms.events` | Kafka topic name for SMS events | Yes |

### Redis Configuration

| Variable Name | Default Value | Description | Required |
|--------------|---------------|-------------|----------|
| `SPRING_DATA_REDIS_HOST` | `redis` | Redis server hostname | Yes |
| `SPRING_DATA_REDIS_PORT` | `6379` | Redis server port | Yes |

### Application-Specific Configuration

| Variable Name | Default Value | Description | Required |
|--------------|---------------|-------------|----------|
| `APP_MOCK_VENDOR_MIN_DELAY_MS` | `100` | Minimum delay (ms) for mock vendor API call | No |
| `APP_MOCK_VENDOR_MAX_DELAY_MS` | `500` | Maximum delay (ms) for mock vendor API call | No |
| `APP_MOCK_VENDOR_FAILURE_RATE` | `0.3` | Probability (0.0-1.0) of vendor API returning failure | No |

---

## Go SMS Store Service

### Server Configuration

| Variable Name | Default Value | Description | Required |
|--------------|---------------|-------------|----------|
| `SERVER_PORT` | `8090` | HTTP server port for the REST API | No |
| `LOG_LEVEL` | `INFO` | Logging level (DEBUG, INFO, WARN, ERROR) | No |

### MongoDB Configuration

| Variable Name | Default Value | Description | Required |
|--------------|---------------|-------------|----------|
| `MONGO_URI` | `mongodb://smsapp:smsapp123@mongodb:27017/sms_store?authSource=sms_store` | Full MongoDB connection URI | Yes |
| `MONGO_DATABASE` | `sms_store` | MongoDB database name | Yes |

**Alternative MongoDB Configuration (if MONGO_URI not provided):**

| Variable Name | Default Value | Description | Required |
|--------------|---------------|-------------|----------|
| `MONGO_HOST` | `mongodb` | MongoDB server hostname | Yes* |
| `MONGO_PORT` | `27017` | MongoDB server port | Yes* |
| `MONGO_APP_USER` | `smsapp` | MongoDB application username | Yes* |
| `MONGO_APP_PASSWORD` | `smsapp123` | MongoDB application password | Yes* |

*Required only if `MONGO_URI` is not provided

### Kafka Configuration

| Variable Name | Default Value | Description | Required |
|--------------|---------------|-------------|----------|
| `KAFKA_BROKERS` | `kafka:9092` | Comma-separated list of Kafka broker addresses | Yes |
| `KAFKA_TOPIC` | `sms.events` | Kafka topic name to consume SMS events from | Yes |
| `KAFKA_GROUP_ID` | `sms-store-consumer-group` | Consumer group ID for Kafka consumer coordination | Yes |

---

## Infrastructure Services

### Zookeeper

| Variable Name | Default Value | Description | Required |
|--------------|---------------|-------------|----------|
| `ZOOKEEPER_CLIENT_PORT` | `2181` | Port for ZooKeeper client connections | No |
| `ZOOKEEPER_TICK_TIME` | `2000` | ZooKeeper tick time in milliseconds | No |
| `ZOOKEEPER_MAX_CLIENT_CNXNS` | `60` | Maximum number of client connections | No |

### Kafka

| Variable Name | Default Value | Description | Required |
|--------------|---------------|-------------|----------|
| `KAFKA_BROKER_ID` | `1` | Unique broker ID for this Kafka instance | No |
| `KAFKA_ZOOKEEPER_CONNECT` | `zookeeper:2181` | ZooKeeper connection string | Yes |
| `KAFKA_ADVERTISED_LISTENERS` | `PLAINTEXT://kafka:9092,PLAINTEXT_HOST://localhost:29092` | Advertised listeners for clients | Yes |
| `KAFKA_LISTENER_SECURITY_PROTOCOL_MAP` | `PLAINTEXT:PLAINTEXT,PLAINTEXT_HOST:PLAINTEXT` | Mapping of listener names to security protocols | No |
| `KAFKA_INTER_BROKER_LISTENER_NAME` | `PLAINTEXT` | Listener name for inter-broker communication | No |
| `KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR` | `1` | Replication factor for offsets topic | No |
| `KAFKA_AUTO_CREATE_TOPICS_ENABLE` | `true` | Auto-create topics on first use | No |
| `KAFKA_LOG_RETENTION_HOURS` | `168` | Log retention period in hours (default: 7 days) | No |
| `KAFKA_LOG_SEGMENT_BYTES` | `1073741824` | Maximum size of a single log segment file | No |
| `KAFKA_NUM_PARTITIONS` | `3` | Default number of partitions for new topics | No |
| `KAFKA_DEFAULT_REPLICATION_FACTOR` | `1` | Default replication factor for topics | No |

### MongoDB

| Variable Name | Default Value | Description | Required |
|--------------|---------------|-------------|----------|
| `MONGO_INITDB_ROOT_USERNAME` | `admin` | MongoDB root username for initial setup | Yes |
| `MONGO_INITDB_ROOT_PASSWORD` | `admin123` | MongoDB root password for initial setup | Yes |
| `MONGO_INITDB_DATABASE` | `sms_store` | Database to be created on initialization | Yes |
| `MONGO_APP_USER` | `smsapp` | Application user to be created | Yes |
| `MONGO_APP_PASSWORD` | `smsapp123` | Application user password | Yes |

### Redis

Redis service uses default configuration with no custom environment variables.

---

## Docker Compose Overrides

These variables can be set in a `.env` file at the project root to override defaults:

```env
# Port Mappings
JAVA_SERVICE_PORT=8080
GO_SERVICE_PORT=8090
KAFKA_PORT=9092
KAFKA_EXTERNAL_PORT=29092
MONGO_PORT=27017
REDIS_PORT=6379
ZOOKEEPER_PORT=2181

# Kafka Configuration
KAFKA_ADVERTISED_HOST=kafka
KAFKA_BROKER_ID=1
KAFKA_TOPIC=sms.events
KAFKA_GROUP_ID=sms-store-consumer-group

# MongoDB Configuration
MONGO_ROOT_USERNAME=admin
MONGO_ROOT_PASSWORD=admin123
MONGO_DATABASE=sms_store
MONGO_APP_USER=smsapp
MONGO_APP_PASSWORD=smsapp123

# Logging
LOG_LEVEL=INFO
SPRING_PROFILES_ACTIVE=docker
```

---

## Security Notes

⚠️ **WARNING**: The default passwords shown in this document are for development purposes only.

For production deployments:

1. **Change all default passwords** to strong, randomly generated values
2. **Use secrets management** (Docker Secrets, Kubernetes Secrets, HashiCorp Vault)
3. **Never commit** `.env` files with real credentials to version control
4. **Rotate credentials** regularly
5. **Use environment-specific** configurations for different deployment stages
6. **Enable TLS/SSL** for all inter-service communication in production
7. **Implement proper** network segmentation and firewall rules

---

## Example .env File

Create a `.env` file in the project root:

```env
# Custom port configuration
JAVA_SERVICE_PORT=8080
GO_SERVICE_PORT=8090

# Production Kafka settings
KAFKA_ADVERTISED_HOST=kafka.production.local
KAFKA_LOG_RETENTION_HOURS=336  # 14 days

# Production MongoDB credentials (use secrets in real prod!)
MONGO_ROOT_USERNAME=admin
MONGO_ROOT_PASSWORD=<secure-password>
MONGO_APP_USER=smsapp
MONGO_APP_PASSWORD=<secure-password>

# Logging for production
LOG_LEVEL=WARN
SPRING_PROFILES_ACTIVE=prod
```

---

## Troubleshooting

### Java Service Can't Connect to Kafka

Check these variables:
- `SPRING_KAFKA_BOOTSTRAP_SERVERS` - Must match Kafka service name/host
- Verify Kafka is running: `docker-compose ps kafka`

### Go Service Can't Connect to MongoDB

Check these variables:
- `MONGO_URI` - Verify connection string format
- `MONGO_APP_USER` and `MONGO_APP_PASSWORD` - Must match MongoDB initialization
- Verify MongoDB is running: `docker-compose ps mongodb`

### Services Can't Communicate

Check these variables:
- `KAFKA_ADVERTISED_HOST` - Must be resolvable by consumers
- Verify Docker network: `docker network inspect polyglot-network`

### Environment Variable Not Taking Effect

1. Ensure variable is defined in `docker-compose.yml` under service's `environment` section
2. Rebuild the service: `docker-compose build <service-name>`
3. Restart with new environment: `docker-compose up -d <service-name>`
4. Verify with: `docker exec <container-name> env | grep <VARIABLE_NAME>`

---

**Last Updated**: December 26, 2025
