# Build stage
FROM golang:1.24.0-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN go build -o main ./cmd/app

# Build the loader
RUN go build -o loader ./cmd/loader

# Run stage
FROM alpine:latest

WORKDIR /app

# Copy the binaries from builder
COPY --from=builder /app/main .
COPY --from=builder /app/loader .

# Copy web assets
COPY --from=builder /app/web ./web

# Copy datasets
COPY datasets ./datasets

# Expose port
EXPOSE 3000

# Run the application
CMD ["./main"]
