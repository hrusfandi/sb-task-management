# Multi-stage build for optimization
# Stage 1: Build stage
FROM golang:1.24-alpine AS builder

# Install dependencies
RUN apk add --no-cache git gcc musl-dev

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application with optimization flags
RUN CGO_ENABLED=0 GOOS=linux go build -mod=readonly -a -installsuffix cgo -ldflags="-w -s" -o main .

# Stage 2: Runtime stage
FROM alpine:latest

# Install ca-certificates for HTTPS, timezone data, and netcat for health checks
RUN apk --no-cache add ca-certificates tzdata netcat-openbsd curl

# Install migrate tool
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.0/migrate.linux-arm64.tar.gz | tar xvz -C /usr/local/bin/ && \
    chmod +x /usr/local/bin/migrate

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/main .

# Copy migration files
COPY --from=builder /app/migrations ./migrations

# Copy entrypoint script
COPY docker-entrypoint.sh /docker-entrypoint.sh
RUN chmod +x /docker-entrypoint.sh

# Expose port
EXPOSE 8080

# Set entrypoint
ENTRYPOINT ["/docker-entrypoint.sh"]

# Run the binary
CMD ["./main"]
