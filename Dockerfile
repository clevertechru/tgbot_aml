FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install git and build dependencies
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bot ./cmd/bot

# Final stage
FROM alpine:latest

WORKDIR /app

# Install ca-certificates for HTTPS and bash for debugging
RUN apk --no-cache add ca-certificates tzdata bash

# Copy the binary from builder
COPY --from=builder /app/bot .
COPY --from=builder /app/config ./config

# Create volume for logs
VOLUME ["/app/logs"]

# Create entrypoint script
COPY <<EOF /app/entrypoint.sh
#!/bin/bash
echo "Starting bot..."
echo "Current directory: $(pwd)"
echo "Config directory contents:"
ls -la /app/config
echo "Environment variables:"
env | grep -E 'TELEGRAM|AML'
echo "Running bot..."
./bot
EOF

RUN chmod +x /app/entrypoint.sh

# Run the application
CMD ["/app/entrypoint.sh"] 