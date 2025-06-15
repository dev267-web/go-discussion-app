FROM golang:1.24 as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Change this path if your main.go is in a subfolder
RUN go build -o main ./cmd/server

FROM debian:bookworm-slim
WORKDIR /app
COPY --from=builder /app/main .

EXPOSE 8080
CMD ["./main"]
