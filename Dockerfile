# Build stage
FROM golang:1.24.3-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git make

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main cmd/api/main.go

# Build migration tool
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o migrate cmd/migrate/migrate.go

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/main .
COPY --from=builder /app/migrate .

# Copy migrations folder
COPY --from=builder /app/internal/database/migrations ./internal/database/migrations

# Expose port
EXPOSE 5000

# Create startup script
RUN echo '#!/bin/sh' > /root/start.sh && \
    echo 'echo "Running database migrations..."' >> /root/start.sh && \
    echo './migrate up' >> /root/start.sh && \
    echo 'echo "Starting application..."' >> /root/start.sh && \
    echo './main' >> /root/start.sh && \
    chmod +x /root/start.sh

# Run the application
CMD ["/root/start.sh"]
