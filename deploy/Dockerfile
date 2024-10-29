FROM golang:1.22.5-alpine AS builder

# Копируем все файлы проекта
COPY .. /github.com/Dnlbb/chat/
WORKDIR /github.com/Dnlbb/chat/

# Скачиваем зависимости
RUN go mod download

# Собираем бинарник
RUN go build -o ./bin/chat-server cmd/main.go

# Финальный этап
FROM alpine:latest

WORKDIR /root/cmd

# Копируем бинарник auth-server из builder стадии
COPY --from=builder /github.com/Dnlbb/chat/bin/chat-server ./cmd/
COPY --from=builder /github.com/Dnlbb/chat/chat.env ./postgres/.env



CMD ["sh", "-c", "cd cmd && ./chat-server"]
