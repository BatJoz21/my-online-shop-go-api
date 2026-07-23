# --- Stage 1: build the binary ---
FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o server main.go

# --- Stage 2: run the binary, nothing else ---
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/server/ .

EXPOSE 8080

CMD ["./server"]