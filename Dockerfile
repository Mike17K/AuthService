# Dockerfile for production build
FROM golang:1.22 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -o auth-service .

# Run the built binary
# CMD ["./auth-service"]


# Start a new stage from scratch
FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/auth-service .
COPY --from=builder /app/.env ./.env

# Run the built binary
CMD ["./auth-service"]