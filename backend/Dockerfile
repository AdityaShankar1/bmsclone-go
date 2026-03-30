# Build stage
FROM golang:1.26-alpine AS builder

WORKDIR /app

# Install dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Final stage
FROM alpine:latest

WORKDIR /root/

# Copy binary
COPY --from=builder /app/main .

# Expose port
EXPOSE 8080

# Environment variables (defaults)
ENV DATABASE_URL="host=db user=aditya password=password123 dbname=bookmyshow port=5432 sslmode=disable"
ENV REDIS_URL="redis:6379"

CMD ["./main"]
