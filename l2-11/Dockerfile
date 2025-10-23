# 1️⃣ Стадия сборки
FROM golang:1.24.5 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o server ./main.go

# 2️⃣ Стадия запуска
FROM alpine:3.19

WORKDIR /app

COPY --from=builder /app/server .
COPY internal/db/migrations ./internal/db/migrations

EXPOSE 8080

CMD ["./server"]
