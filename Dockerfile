# Stage 1: Build the Go app
FROM golang:1.21 AS builder

# Set Go proxy to avoid module download failures
ENV GOPROXY=https://proxy.golang.org,direct

WORKDIR /app

# Copy go.mod and go.sum first for dependency caching
COPY go.mod ./
COPY go.sum ./
RUN go mod download

# Now copy the entire source code
COPY . ./

# Build the Go binary
RUN go build -o server ./cmd/main.go

# Stage 2: Minimal image to run the app
FROM gcr.io/distroless/base-debian10

WORKDIR /app

# Copy the compiled binary from builder
COPY --from=builder /app/server .

# Expose port 8080 (or the one you use)
EXPOSE 8080

# Run the app
CMD ["./server"]
