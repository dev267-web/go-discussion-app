# Use an official Go image with the required version
FROM golang:1.24 as builder

WORKDIR /app

# Copy go.mod and go.sum first to cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire source code
COPY . .

# Build the Go app
RUN go build -o main .

# --------- Optional: final small image for deployment ---------
FROM debian:bookworm-slim

WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /app/main .

# Expose a port if your app uses one (adjust as needed)
EXPOSE 8080

# Run the app
CMD ["./main"]
