# Build stage
FROM golang:1.23.0-alpine AS builder

WORKDIR /app
COPY . .

RUN go mod download
RUN go build -o main ./main.go

# Run stage
FROM alpine:3.14

WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/internal/config ./config

EXPOSE 8080
CMD ["./main"]