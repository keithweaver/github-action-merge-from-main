# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /build

# Copy go module files first for better caching
COPY go.mod .

# Copy source files
COPY *.go .

# Build the Go binary
RUN go build -o pr-merger .

# Runtime stage
FROM alpine:latest

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /build/pr-merger /app/pr-merger

# Copy entrypoint script
COPY entrypoint.sh /app/entrypoint.sh
RUN chmod +x /app/entrypoint.sh

ENTRYPOINT ["/app/entrypoint.sh"]
