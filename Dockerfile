# Build stage
FROM golang:1.23-alpine AS builder

# Install build dependencies
RUN apk add --no-cache \
    build-base \
    vips-dev \
    pkgconfig

WORKDIR /app

# Copy and download dependencies first (for better caching)
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN go build -o frame ./cmd/api

# Runtime stage
FROM alpine:latest

# Install runtime dependencies for libvips
RUN apk add --no-cache \
    vips \
    ca-certificates

# Create a non-root user to run the application
RUN adduser -D -g '' appuser

WORKDIR /app

# Copy the binary from builder stage
COPY --from=builder /app/frame .

# Use non-root user for security
USER appuser

# Expose the API port
EXPOSE 8080

# Run the application
CMD ["/app/frame"]