FROM golang:1.24-alpine AS builder

# Set working directory
WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git gcc musl-dev

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=1 GOOS=linux go build -a -ldflags "-linkmode external -extldflags '-static' -s -w" -o accio ./cmd/accio

# Create a minimal image
FROM alpine:latest

# Install runtime dependencies
RUN apk add --no-cache ca-certificates tzdata sqlite

# Set working directory
WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/accio /app/accio

# Copy web assets
COPY web /app/web

# Create cache directory
RUN mkdir -p /app/cache/images

# Set environment variables
ENV PROFILE_IMAGE_CACHE_DIR=/app/cache/images

# Expose port for web UI
EXPOSE 8080

# Run the application
ENTRYPOINT ["/app/accio"]
CMD ["--web", "--use-database"]