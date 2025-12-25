# Multi-stage build for Go SMS Store Service
# Stage 1: Build with Go
FROM golang:1.25-alpine AS build

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git

# Copy go mod files and download dependencies (cached layer)
COPY go.mod go.sum* ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o sms-store .

# Stage 2: Runtime with minimal Alpine
FROM alpine:latest

WORKDIR /app

# Install ca-certificates for HTTPS and wget for health check
RUN apk --no-cache add ca-certificates wget

# Create non-root user for security
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
USER appuser:appgroup

# Copy the built binary from build stage
COPY --from=build /app/sms-store .

# Expose application port
EXPOSE 8090

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=40s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8090/health || exit 1

# Run the application
ENTRYPOINT ["./sms-store"]
