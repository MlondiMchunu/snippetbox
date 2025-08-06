# Build stage
FROM golang:1.24.2-alpine AS builder

# Install dependencies
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o snippetbox ./cmd/web

# Runtime stage
FROM alpine:latest

# Create appuser
RUN adduser -D -g '' appuser

# Set working directory
WORKDIR /home/appuser

# Copy certificates
COPY --from=builder /app/certificates ./certificates
COPY --from=builder /app/tls ./tls

# Copy static files
COPY --from=builder /app/ui/static ./ui/static
COPY --from=builder /app/ui/html ./ui/html

# Copy binary
COPY --from=builder --chown=appuser /app/snippetbox .

# Expose ports
EXPOSE 4000
EXPOSE 443

# Run as non-root user
USER appuser

# Start the application
CMD ["./snippetbox"]