FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

# Download all deps
RUN go mod download -x

COPY . .

# Build the application
RUN go build -o main ./cmd/gollet/main.go

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/main .

COPY .env .

# Expose port 8080
EXPOSE 8080

# Run the application
CMD ["./main"]