FROM golang:1.23-alpine AS builder

WORKDIR /app

# Копируем go.mod и go.sum для кэширования зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем весь пакет main-service
RUN CGO_ENABLED=0 go build -o main-service ./cmd/main-service

# Финальный образ
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

# Копируем только исполняемый файл
COPY --from=builder /app/main-service .

# Создаем папку для логов
RUN mkdir -p logs

# Expose оба порта
EXPOSE 37545 33005

# Запускаем с прямым указанием на docker конфиг
CMD ["./main-service", "-config", "configs/config.docker.yaml"] 